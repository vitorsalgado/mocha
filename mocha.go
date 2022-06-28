package mocha

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/params"
)

type (
	configT interface{ *Config | *ConfigBuilder }

	Mocha struct {
		Server      *httptest.Server
		mockStorage mock.Storage
		context     context.Context
		params      params.Params
	}
)

func New[C configT](config C) *Mocha {
	var parsedConfig *Config
	switch conf := any(config).(type) {
	case *ConfigBuilder:
		parsedConfig = conf.Build()
	case *Config:
		parsedConfig = conf
	}

	if parsedConfig == nil {
		parsedConfig = Setup().Build()
	}

	mockStorage := mock.NewStorage()
	parsers := make([]BodyParser, 0)
	parsers = append(parsers, parsedConfig.BodyParsers...)
	parsers = append(parsers, &JSONBodyParser{}, &FormURLEncodedParser{})
	params := params.New()
	params.Set(BuiltIntExtraScenario, NewScenarioStore())

	return &Mocha{
		Server:      httptest.NewUnstartedServer(newHandler(mockStorage, parsers, params)),
		mockStorage: mockStorage,
		context:     parsedConfig.Context,
		params:      params}
}

func ConfigureForTest[C configT](t *testing.T, options C) *Mocha {
	m := New(options)
	t.Cleanup(func() { m.Close() })
	return m
}

func ForTest(t *testing.T) *Mocha {
	return ConfigureForTest(t, Setup())
}

func (m *Mocha) Start() ServerInfo {
	m.Server.Start()
	return ServerInfo{URL: m.Server.URL}
}

func (m *Mocha) StartTLS() ServerInfo {
	m.Server.StartTLS()
	return ServerInfo{URL: m.Server.URL}
}

func (m *Mocha) Mock(builders ...*MockBuilder) *Scoped {
	size := len(builders)
	added := make([]*mock.Mock, size)

	for i, b := range builders {
		newMock := b.Build()
		m.mockStorage.Save(newMock)
		added[i] = newMock
	}

	return Scope(m.mockStorage, added)
}

func (m *Mocha) Parameters() *params.Params {
	return &m.params
}

func (m *Mocha) Close() {
	m.Server.Close()
}
