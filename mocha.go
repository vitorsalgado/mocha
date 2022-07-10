package mocha

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/expect/scenario"
	"github.com/vitorsalgado/mocha/internal/middleware"
	"github.com/vitorsalgado/mocha/internal/parameters"
)

type (
	configT interface{ *Config | *Configurer }

	// Mocha is the base for the mock server.
	Mocha struct {
		server  *httptest.Server
		storage core.Storage
		context context.Context
		params  parameters.Params
		t       core.T
	}
)

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a Configurer implementation.
func New[C configT](t core.T, config C) *Mocha {
	var cfg *Config
	switch conf := any(config).(type) {
	case *Configurer:
		cfg = conf.Build()
	case *Config:
		cfg = conf
	}

	if cfg == nil {
		cfg = Configure().Build()
	}

	storage := core.NewStorage()

	parsers := make([]RequestBodyParser, 0)
	parsers = append(parsers, cfg.BodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &bytesParser{})

	params := parameters.New()
	params.Set(scenario.BuiltInParamStore, scenario.NewStore())

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, middleware.Recover)
	middlewares = append(middlewares, cfg.Middlewares...)

	handler := middleware.Compose(middlewares...).Root(newHandler(storage, parsers, params, t))

	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true

	m := &Mocha{
		server:  server,
		storage: storage,
		context: cfg.Context,
		params:  params,
		t:       t}

	t.Cleanup(func() { m.Close() })

	return m
}

// ForTest creates a new Mocha mock server with default configurations.
// It closes the mock server after the tests finishes, using the testing.T cleanup feature.
func ForTest(t core.T) *Mocha {
	return New(t, Configure())
}

// Start starts the mock server.
func (m *Mocha) Start() ServerInfo {
	m.server.Start()
	return ServerInfo{URL: m.server.URL}
}

// StartTLS starts TLS from a server.
func (m *Mocha) StartTLS() ServerInfo {
	m.server.StartTLS()
	return ServerInfo{URL: m.server.URL}
}

// Mock adds one or multiple HTTP request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checking if they were called or not.
// The returned Scoped is useful for tests.
//
// Example:
// 	scoped := m.Mock(
// 		Get(to.URLPath("/test")).
// 			Header("test", to.EqualTo("hello")).
// 			Query("filter", to.EqualTo("all")).
// 			Reply(reply.
// 				Created().
// 				BodyString("hello world")))
//
//	assert.True(t, scoped.Called())
func (m *Mocha) Mock(builders ...*MockBuilder) Scoped {
	size := len(builders)
	added := make([]*core.Mock, size)

	for i, b := range builders {
		newMock := b.Build()
		m.storage.Save(newMock)
		added[i] = newMock
	}

	return scope(m.storage, added)
}

// Parameters allows managing custom parameters that will be available inside matchers.
func (m *Mocha) Parameters() parameters.Params {
	return m.params
}

// URL returns the base URL string for the mock server.
func (m *Mocha) URL() string {
	return m.server.URL
}

// Close closes the mock server.
func (m *Mocha) Close() {
	m.server.Close()
}
