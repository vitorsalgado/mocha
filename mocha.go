package mocha

import (
	"context"
	"net/http"
	"sync"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/feat/cors"
	"github.com/vitorsalgado/mocha/feat/scenario"
	"github.com/vitorsalgado/mocha/internal/middleware"
	"github.com/vitorsalgado/mocha/internal/parameters"
)

type (
	// Mocha is the base for the mock server.
	Mocha struct {
		server  Server
		storage core.Storage
		context context.Context
		cancel  context.CancelFunc
		params  parameters.Params
		scopes  []*Scoped
		mu      *sync.Mutex
		t       core.T
	}

	// Cleanable allows marking mocha instance to be closed on test cleanup.
	Cleanable interface {
		Cleanup(func())
	}
)

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a Configurer implementation.
func New(t core.T, config ...Config) *Mocha {
	t.Helper()

	cfg := configDefault
	if len(config) > 0 {
		cfg = config[0]
	}

	storage := core.NewStorage()

	parsers := make([]RequestBodyParser, 0)
	parsers = append(parsers, cfg.BodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &bytesParser{})

	params := parameters.New()
	params.Set(scenario.BuiltInParamStore, scenario.NewStore())

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, middleware.Recover)

	if cfg.corsEnabled {
		middlewares = append(middlewares, cors.New(cfg.CORS))
	}

	middlewares = append(middlewares, cfg.Middlewares...)

	handler := middleware.
		Compose(middlewares...).
		Root(newHandler(storage, parsers, params, t))

	server := cfg.Server

	if server == nil {
		server = newServer()
	}

	err := server.Configure(cfg, handler)
	if err != nil {
		t.Errorf("failed to configure mock server. reason: %v", err)
		t.FailNow()
	}

	parent := cfg.Context
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	m := &Mocha{
		server:  server,
		storage: storage,
		context: ctx,
		cancel:  cancel,
		params:  params,
		scopes:  make([]*Scoped, 0),
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

// NewSimple creates a new Mocha mock server with default configurations.
// It closes the mock server after the tests finishes, using the testing.T cleanup feature.
func NewSimple() *Mocha {
	return New(&noop{})
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
func (m *Mocha) Mock(builders ...*MockBuilder) *Scoped {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := len(builders)
	added := make([]*core.Mock, size)

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
func (m *Mocha) Parameters() parameters.Params {
	return m.params
}

// URL returns the base URL string for the mock server.
func (m *Mocha) URL() string {
	return m.server.Info().URL
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
func (m *Mocha) AssertCalled(t core.T) bool {
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
func (m *Mocha) AssertNotCalled(t core.T) bool {
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
func (m *Mocha) AssertHits(t core.T, expected int) bool {
	t.Helper()

	hits := m.Hits()

	if hits < expected {
		t.Errorf("\nexpected %d request hits. got %d", expected, hits)
		return false
	}

	return true
}
