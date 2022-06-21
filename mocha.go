package mocha

import (
	"context"
	"net/http/httptest"
	"testing"
)

type (
	Mocha struct {
		Server    *httptest.Server
		mockstore MockStore
		context   context.Context
	}
)

func New(context context.Context) *Mocha {
	mockstore := NewMockStore()
	parsers := make([]BodyParser, 0)
	parsers = append(parsers, &JSONBodyParser{}, &FormURLEncodedParser{})
	extras := NewExtras()
	extras.Set(BuiltIntExtraScenario, NewScenarioStore())

	return &Mocha{
		Server:    httptest.NewUnstartedServer(newHandler(mockstore, parsers, extras)),
		mockstore: mockstore,
		context:   context}
}

func ForTestWithContext(t *testing.T, context context.Context) *Mocha {
	m := New(context)
	t.Cleanup(func() { m.Close() })
	return m
}

func ForTest(t *testing.T) *Mocha {
	return ForTestWithContext(t, context.Background())
}

func (m *Mocha) Start() ServerInfo {
	m.Server.Start()

	go func() {
		<-m.context.Done()
		m.Close()
	}()

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
