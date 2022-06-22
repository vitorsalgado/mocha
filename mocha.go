package mocha

import (
	"context"
	"net/http/httptest"
	"testing"
)

type (
	ConfT interface{ *Config | *ConfigBuilder }

	Mocha struct {
		Server    *httptest.Server
		mockstore MockStore
		context   context.Context
	}
)

func New[C ConfT](options C) *Mocha {
	var opts *Config
	switch conf := any(options).(type) {
	case *ConfigBuilder:
		opts = conf.Build()
	case *Config:
		opts = conf
	}

	if opts == nil {
		opts = Options().Build()
	}

	mockstore := NewMockStore()
	parsers := make([]BodyParser, 0)
	parsers = append(parsers, &JSONBodyParser{}, &FormURLEncodedParser{})
	extras := NewExtras()
	extras.Set(BuiltIntExtraScenario, NewScenarioStore())

	return &Mocha{
		Server:    httptest.NewUnstartedServer(newHandler(mockstore, parsers, extras)),
		mockstore: mockstore,
		context:   opts.Context}
}

func ConfigureForTest[C ConfT](t *testing.T, options C) *Mocha {
	m := New(options)
	t.Cleanup(func() { m.Close() })
	return m
}

func ForTest(t *testing.T) *Mocha {
	return ConfigureForTest(t, Options())
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
	added := make([]int32, 0)

	for _, b := range builders {
		mock := b.Build()

		m.mockstore.Save(mock)
		added = append(added, mock.ID)
	}

	return NewScoped(m.mockstore, added)
}

func (m *Mocha) Close() {
	m.Server.Close()
}
