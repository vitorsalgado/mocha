package mocha

import (
	"context"
	"net/http"
	"reflect"
	"sync"

	"github.com/vitorsalgado/mocha/v3/cors"
	"github.com/vitorsalgado/mocha/v3/hooks"
	"github.com/vitorsalgado/mocha/v3/internal/mid"
	"github.com/vitorsalgado/mocha/v3/internal/mid/recover"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type (
	// Mocha is the base for the mock server.
	Mocha struct {
		Config *Config
		T      TestingT

		server  Server
		storage storage
		ctx     context.Context
		cancel  context.CancelFunc
		params  reply.Params
		hooks   *hooks.Hooks
		scopes  []*Scoped
		mu      sync.Mutex
	}

	// Cleanable allows marking mocha instance to be closed on test cleanup.
	Cleanable interface {
		Cleanup(func())
	}
)

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a Configurer implementation.
func New(t TestingT, config ...*Config) *Mocha {
	cfg := _configDefault
	if len(config) > 0 {
		cfg = config[0]
	}

	ctx, cancel := context.WithCancel(context.Background())
	mockStorage := newStorage()
	parsers := make([]RequestBodyParser, 0)
	parsers = append(parsers, cfg.BodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &bytesParser{})

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, recover.Recover)

	hook := hooks.New()

	if cfg.LogLevel > LogSilently {
		h := hooks.NewInternalEvents(t)

		hook.Subscribe(hooks.HookOnRequest, h.OnRequest)
		hook.Subscribe(hooks.HookOnRequestMatched, h.OnRequestMatched)
		hook.Subscribe(hooks.HookOnRequestNotMatched, h.OnRequestNotMatched)
		hook.Subscribe(hooks.HookOnError, h.OnError)
	}

	if cfg.corsEnabled {
		middlewares = append(middlewares, cors.New(cfg.CORS))
	}

	middlewares = append(middlewares, cfg.Middlewares...)
	p := reply.Parameters()
	handler := mid.
		Compose(middlewares...).
		Root(newHandler(mockStorage, parsers, p, hook, t))

	if cfg.Handler != nil {
		handler = cfg.Handler(handler)
	}

	server := cfg.Server

	if server == nil {
		server = newServer()
	}

	err := server.Configure(cfg, handler)
	if err != nil {
		t.Errorf("failed to configure mock server. reason=%v", err)
		t.FailNow()
	}

	m := &Mocha{
		Config: cfg,

		server:  server,
		storage: mockStorage,
		ctx:     ctx,
		cancel:  cancel,
		params:  p,
		scopes:  make([]*Scoped, 0),
		hooks:   hook,
		T:       t}

	return m
}

// NewBasic creates a new Mocha mock server with default configurations.
func NewBasic() *Mocha {
	return New(NewConsoleNotifier())
}

// Start starts the mock server.
func (m *Mocha) Start() ServerInfo {
	info, err := m.server.Start()
	if err != nil {
		m.T.Errorf("failed to start mock server. reason=%v", err)
		m.T.FailNow()
	}

	m.hooks.Start(m.ctx)

	return info
}

// StartTLS starts TLS from a server.
func (m *Mocha) StartTLS() ServerInfo {
	info, err := m.server.StartTLS()
	if err != nil {
		m.T.Errorf("failed to start a TLS mock server. reason=%v", err)
		m.T.FailNow()
	}

	m.hooks.Start(m.ctx)

	return info
}

// AddMocks adds one or multiple request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checking if they were called or not.
// The returned Scoped is useful for tests.
//
// Usage:
//
//	scoped := m.AddMocks(
//		Get(matcher.URLPath("/test")).
//			Header("test", matcher.Equal("hello")).
//			Query("filter", matcher.Equal("all")).
//			Reply(reply.Created().BodyString("hello world")))
//
//	assert.True(T, scoped.Called())
func (m *Mocha) AddMocks(builders ...*MockBuilder) *Scoped {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := len(builders)
	added := make([]*Mock, size)

	for i, b := range builders {
		nm := b.Build(&Deps{ScenarioStore: matcher.NewScenarioStorage()})
		m.storage.Save(nm)
		added[i] = nm
	}

	scoped := scope(m.storage, added)
	m.scopes = append(m.scopes, scoped)

	return scoped
}

// Parameters allows managing custom parameters that will be available inside matchers.
func (m *Mocha) Parameters() reply.Params {
	return m.params
}

// URL returns the base URL string for the mock server.
func (m *Mocha) URL() string {
	return m.server.Info().URL
}

// Subscribe add a new event listener.
func (m *Mocha) Subscribe(evt reflect.Type, fn func(payload any)) {
	m.hooks.Subscribe(evt, fn)
}

// Close closes the mock server.
func (m *Mocha) Close() {
	m.cancel()

	err := m.server.Close()
	if err != nil {
		m.T.Logf(err.Error())
	}
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
	t.Cleanup(func() { m.Close() })
	return m
}

// AssertCalled asserts that all mocks associated with this instance were called at least once.
func (m *Mocha) AssertCalled(t TestingT) bool {
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
func (m *Mocha) AssertNotCalled(t TestingT) bool {
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
func (m *Mocha) AssertHits(t TestingT, expected int) bool {
	t.Helper()

	hits := m.Hits()

	if hits < expected {
		t.Errorf("\nexpected %d request hits. got %d", expected, hits)
		return false
	}

	return true
}
