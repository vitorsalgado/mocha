package mocha

import (
	"context"
	"net/http"
	"sync"

	"github.com/vitorsalgado/mocha/v2/cors"
	"github.com/vitorsalgado/mocha/v2/expect"
	"github.com/vitorsalgado/mocha/v2/hooks"
	"github.com/vitorsalgado/mocha/v2/internal/middleware"
	"github.com/vitorsalgado/mocha/v2/internal/middleware/recover"
	"github.com/vitorsalgado/mocha/v2/params"
)

type (
	// Mocha is the base for the mock server.
	Mocha struct {
		server  Server
		storage storage
		context context.Context
		cancel  context.CancelFunc
		params  params.P
		events  *hooks.Emitter
		scopes  []*Scoped
		mu      *sync.Mutex
		t       T
	}

	// Cleanable allows marking mocha instance to be closed on test cleanup.
	Cleanable interface {
		Cleanup(func())
	}
)

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a Configurer implementation.
func New(t T, config ...Config) *Mocha {
	t.Helper()

	cfg := configDefault
	if len(config) > 0 {
		cfg = config[0]
	}

	parent := cfg.Context
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	mockStorage := newStorage()

	parsers := make([]RequestBodyParser, 0)
	parsers = append(parsers, cfg.BodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &bytesParser{})

	p := params.New()
	p.Set(expect.ScenarioBuiltInParamStore, expect.NewScenarioStore())

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, recover.Recover)

	evt := hooks.NewEmitter(ctx)

	if cfg.LogVerbosity == LogVerbose {
		evt.Subscribe(hooks.NewInternalEvents(t))
	}

	if cfg.corsEnabled {
		middlewares = append(middlewares, cors.New(cfg.CORS))
	}

	middlewares = append(middlewares, cfg.Middlewares...)

	handler := middleware.
		Compose(middlewares...).
		Root(newHandler(mockStorage, parsers, p, evt, t))

	server := cfg.Server

	if server == nil {
		server = newServer()
	}

	err := server.Configure(cfg, handler)
	if err != nil {
		t.Errorf("failed to configure mock server. reason: %v", err)
		t.FailNow()
	}

	m := &Mocha{
		server:  server,
		storage: mockStorage,
		context: ctx,
		cancel:  cancel,
		params:  p,
		scopes:  make([]*Scoped, 0),
		events:  evt,
		mu:      &sync.Mutex{},
		t:       t}

	go func() {
		<-ctx.Done()
		e := m.Close()
		if e != nil {
			m.t.Logf("\nerror closing mocha http server. error=%v", e)
		}
	}()

	return m
}

// NewBasic creates a new Mocha mock server with default configurations.
func NewBasic() *Mocha {
	return New(NewConsoleNotifier())
}

// Start starts the mock server.
func (m *Mocha) Start() ServerInfo {
	m.t.Helper()

	info, err := m.server.Start()
	if err != nil {
		m.t.Errorf("failed to start mock server. reason: %v", err)
		m.t.FailNow()
	}

	return info
}

// StartTLS starts TLS from a server.
func (m *Mocha) StartTLS() ServerInfo {
	m.t.Helper()

	info, err := m.server.StartTLS()
	if err != nil {
		m.t.Errorf("failed to start a TLS mock server. reason: %v", err)
		m.t.FailNow()
	}

	return info
}

// AddMocks adds one or multiple request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checking if they were called or not.
// The returned Scoped is useful for tests.
//
// Usage:
//
//	scoped := m.AddMocks(
//		Get(expect.URLPath("/test")).
//			Header("test", expect.ToEqual("hello")).
//			Query("filter", expect.ToEqual("all")).
//			Reply(reply.Created().BodyString("hello world")))
//
//	assert.True(t, scoped.Called())
func (m *Mocha) AddMocks(builders ...*MockBuilder) *Scoped {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := len(builders)
	added := make([]*Mock, size)

	for i, b := range builders {
		newMock := b.Build()
		m.storage.Save(newMock)
		added[i] = newMock
	}

	scoped := scope(m.storage, added)
	m.scopes = append(m.scopes, scoped)

	return scoped
}

// Parameters allows managing custom parameters that will be available inside matchers.
func (m *Mocha) Parameters() params.P {
	return m.params
}

// URL returns the base URL string for the mock server.
func (m *Mocha) URL() string {
	return m.server.Info().URL
}

// Subscribe add a new event listener.
func (m *Mocha) Subscribe(evt hooks.Events) {
	m.events.Subscribe(evt)
}

// Close closes the mock server.
func (m *Mocha) Close() error {
	return m.server.Close()
}

// Hits returns the total request hits.
func (m *Mocha) Hits() int {
	hits := 0

	for _, s := range m.scopes {
		hits += s.Hits()
	}

	return hits
}

// Disable disables all mocks.
func (m *Mocha) Disable() {
	for _, scoped := range m.scopes {
		scoped.Disable()
	}
}

// Enable disables all mocks.
func (m *Mocha) Enable() {
	for _, scoped := range m.scopes {
		scoped.Enable()
	}
}

// CloseOnCleanup adds mocha server Close to the Cleanup.
func (m *Mocha) CloseOnCleanup(t Cleanable) *Mocha {
	closeSrv := func() {
		e := m.Close()
		if e != nil {
			m.t.Logf("\nerror closing mocha http server. error=%v", e)
		}
	}

	t.Cleanup(func() {
		defer m.cancel()
		closeSrv()
	})

	return m
}

// AssertCalled asserts that all mocks associated with this instance were called at least once.
func (m *Mocha) AssertCalled(t T) bool {
	t.Helper()

	result := true

	for i, s := range m.scopes {
		if s.IsPending() {
			t.Logf("\nscope #%d\n", i)
			s.AssertCalled(t)
			result = false
		}
	}

	return result
}

// AssertNotCalled asserts that all mocks associated with this instance were called at least once.
func (m *Mocha) AssertNotCalled(t T) bool {
	t.Helper()

	result := true

	for i, s := range m.scopes {
		if !s.IsPending() {
			t.Logf("\nscope #%d\n", i)
			s.AssertNotCalled(t)
			result = false
		}
	}

	return result
}

// AssertHits asserts that the sum of request hits for mocks
// is equal to the given expected value.
func (m *Mocha) AssertHits(t T, expected int) bool {
	t.Helper()

	hits := m.Hits()

	if hits < expected {
		t.Errorf("\nexpected %d request hits. got %d", expected, hits)
		return false
	}

	return true
}
