package mocha

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/internal/scenario"
	"github.com/vitorsalgado/mocha/mock"
)

type (
	configT interface{ *Config | *Configurer }

	// Mocha is the base for the mock server.
	Mocha struct {
		Server      *httptest.Server
		mockStorage mock.Storage
		context     context.Context
		params      params.Params
	}
)

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a Configurer implementation.
func New[C configT](config C) *Mocha {
	var parsedConfig *Config
	switch conf := any(config).(type) {
	case *Configurer:
		parsedConfig = conf.Build()
	case *Config:
		parsedConfig = conf
	}

	if parsedConfig == nil {
		parsedConfig = Configure().Build()
	}

	mockStorage := mock.NewStorage()
	parsers := make([]RequestBodyParser, 0)
	parsers = append(parsers, parsedConfig.BodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &bytesParser{})
	parameters := params.New()
	parameters.Set(scenario.BuiltInParamStore, scenario.NewStore())

	server := httptest.NewUnstartedServer(newHandler(mockStorage, parsers, parameters))
	server.EnableHTTP2 = true

	return &Mocha{
		Server:      server,
		mockStorage: mockStorage,
		context:     parsedConfig.Context,
		params:      parameters}
}

// ConfigureForTest creates a new Mocha mock server with the given configurations.
// It closes the mock server after the tests finishes, using the testing.T cleanup feature.
// Parameter config accepts a Config or a Configurer implementation.
func ConfigureForTest[C configT](t *testing.T, options C) *Mocha {
	m := New(options)
	t.Cleanup(func() { m.Close() })
	return m
}

// ForTest creates a new Mocha mock server with default configurations.
// It closes the mock server after the tests finishes, using the testing.T cleanup feature.
func ForTest(t *testing.T) *Mocha {
	return ConfigureForTest(t, Configure())
}

// Start starts the mock server.
func (m *Mocha) Start() ServerInfo {
	m.Server.Start()
	return ServerInfo{URL: m.Server.URL}
}

// StartTLS starts TLS from a server.
func (m *Mocha) StartTLS() ServerInfo {
	m.Server.StartTLS()
	return ServerInfo{URL: m.Server.URL}
}

// Mock adds one or multiple HTTP request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checking if they were called or not.
// The returned Scoped is useful for tests.
//
// Example:
// 	scoped := m.Mock(
// 		Get(matchers.URLPath("/test")).
// 			Header("test", matchers.EqualTo("hello")).
// 			Query("filter", matchers.EqualTo("all")).
// 			Reply(
// 				reply.
// 					Created().
// 					BodyString("hello world")))
//
//	assert.True(t, scoped.IsDone())
func (m *Mocha) Mock(builders ...*MockBuilder) Scoped {
	size := len(builders)
	added := make([]*mock.Mock, size)

	for i, b := range builders {
		newMock := b.Build()
		m.mockStorage.Save(newMock)
		added[i] = newMock
	}

	return scope(m.mockStorage, added)
}

// Parameters allows managing custom parameters that will be available inside matchers.
func (m *Mocha) Parameters() params.Params {
	return m.params
}

// Close closes the mock server.
func (m *Mocha) Close() {
	m.Server.Close()
}
