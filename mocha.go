package mocha

import (
	"context"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/templating"
	"net/http/httptest"
	"testing"
)

type (
	configT interface{ *Config | *ConfigBuilder }

	Mocha struct {
		Server         *httptest.Server
		mockStorage    mock.Storage
		context        context.Context
		templateParser templating.Parser
	}
)

func New[C configT](config C) *Mocha {
	var opts *Config
	switch conf := any(config).(type) {
	case *ConfigBuilder:
		opts = conf.Build()
	case *Config:
		opts = conf
	}

	if opts == nil {
		opts = Setup().Build()
	}

	mockStorage := mock.NewMockStorage()
	parsers := make([]BodyParser, len(opts.BodyParsers)+2)
	parsers = append(parsers, opts.BodyParsers...)
	parsers = append(parsers, &JSONBodyParser{}, &FormURLEncodedParser{})
	extras := NewExtras()
	extras.Set(BuiltIntExtraScenario, NewScenarioStore())

	return &Mocha{
		Server:      httptest.NewUnstartedServer(newHandler(mockStorage, parsers, extras)),
		mockStorage: mockStorage,
		context:     opts.Context}
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
	added := make([]int32, len(builders))

	for _, b := range builders {
		newMock := b.Build()

		m.mockStorage.Save(newMock)
		added = append(added, newMock.ID)
	}

	return scoped(m.mockStorage, added)
}

func (m *Mocha) Close() {
	m.Server.Close()
}
