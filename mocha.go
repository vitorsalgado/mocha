package mocha

import (
	"context"
	"net/http"
	"reflect"
	"sync"

	"github.com/vitorsalgado/mocha/v3/internal/mid"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type (
	// Mocha is the base for the mock server.
	Mocha struct {
		Config *Config
		T      TestingT

		server  Server
		storage mockStore
		ctx     context.Context
		cancel  context.CancelFunc
		params  reply.Params
		events  *eventListener
		scopes  []*Scoped
		loaders []Loader
		mu      sync.Mutex
	}

	// TestingT is based on testing.T and allow mocha components to log information and errors.
	TestingT interface {
		Helper()
		Logf(format string, a ...any)
		Errorf(format string, a ...any)
		FailNow()
	}

	// Cleanable allows marking mocha instance to be closed on test cleanup.
	Cleanable interface {
		Cleanup(func())
	}
)

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a Configurer implementation.
func New(t TestingT, config ...*Config) *Mocha {
	conf := _configDefault
	if len(config) > 0 {
		conf = config[0]
	}

	ctx, cancel := context.WithCancel(context.Background())
	store := newStore()
	parsers := make([]RequestBodyParser, 0, len(conf.BodyParsers)+4)
	parsers = append(parsers, conf.BodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &bytesParser{})

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, mid.Recover)

	events := newEvents()

	if conf.LogLevel > LogSilently {
		h := newInternalEvents(t)

		events.Subscribe(EventOnRequest, h.OnRequest)
		events.Subscribe(EventOnRequestMatched, h.OnRequestMatched)
		events.Subscribe(EventOnRequestNotMatched, h.OnRequestNotMatched)
		events.Subscribe(EventOnError, h.OnError)
	}

	if conf.corsEnabled {
		middlewares = append(middlewares, corsMid(conf.CORS))
	}

	middlewares = append(middlewares, conf.Middlewares...)

	p := reply.Parameters()
	if conf.Parameters != nil {
		p = conf.Parameters
	}

	handler := mid.
		Compose(middlewares...).
		Root(newHandler(store, parsers, p, events, t))

	if conf.Handler != nil {
		handler = conf.Handler(handler)
	}

	server := conf.Server

	if server == nil {
		server = newServer()
	}

	err := server.Configure(conf, handler)
	if err != nil {
		t.Errorf("failed to configure mock server. reason=%v", err)
		t.FailNow()
	}

	loaders := make([]Loader, 0)
	loaders = append(loaders, &FileLoader{})

	m := &Mocha{
		Config: conf,

		server:  server,
		storage: store,
		ctx:     ctx,
		cancel:  cancel,
		params:  p,
		scopes:  make([]*Scoped, 0),
		events:  events,
		loaders: loaders,
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

	m.Rebuild()

	m.events.Start(m.ctx)

	return info
}

// StartTLS starts TLS from a server.
func (m *Mocha) StartTLS() ServerInfo {
	info, err := m.server.StartTLS()
	if err != nil {
		m.T.Errorf("failed to start a TLS mock server. reason=%v", err)
		m.T.FailNow()
	}

	m.Rebuild()

	m.events.Start(m.ctx)

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
//			Reply(reply.Created().PlainText("hello world")))
//
//	assert.True(T, scoped.Called())
func (m *Mocha) AddMocks(builders ...Builder) *Scoped {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := len(builders)
	added := make([]*Mock, size)

	for i, b := range builders {
		nm := b.Build()
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

// Context returns internal context.Context.
func (m *Mocha) Context() context.Context {
	return m.ctx
}

// Subscribe add a new event listener.
func (m *Mocha) Subscribe(evt reflect.Type, fn func(payload any)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events.Subscribe(evt, fn)
}

func (m *Mocha) Loader(loader Loader) {
	m.loaders = append(m.loaders, loader)
}

func (m *Mocha) Rebuild() {
	for _, loader := range m.loaders {
		err := loader.Load(m)
		if err != nil {
			panic(err)
		}
	}
}

// Close closes the mock server.
func (m *Mocha) Close() {
	m.cancel()

	err := m.server.Close()
	if err != nil {
		m.T.Logf(err.Error())
	}
}

// CloseOnCleanup adds mocha server Close to the Cleanup.
func (m *Mocha) CloseOnCleanup(t Cleanable) *Mocha {
	t.Cleanup(func() { m.Close() })
	return m
}

// Hits returns the total request hits.
func (m *Mocha) Hits() int {
	hits := 0

	for _, s := range m.scopes {
		hits += s.Hits()
	}

	return hits
}

// Enable enables all mocks.
func (m *Mocha) Enable() {
	for _, scoped := range m.scopes {
		scoped.Enable()
	}
}

// Disable disables all mocks.
func (m *Mocha) Disable() {
	for _, scoped := range m.scopes {
		scoped.Disable()
	}
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
