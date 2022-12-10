package mocha

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/vitorsalgado/mocha/v3/internal/mid"
	"github.com/vitorsalgado/mocha/v3/internal/notifier"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

// StatusNoMockFound describes an HTTP response where no Mock was found.
//
// It uses http.StatusTeapot to reduce the chance of using the same
// expected response from the actual server being mocked.
// Basically, every request that doesn't match against to a Mock will return http.StatusTeapot.
const StatusNoMockFound = http.StatusTeapot

// Mocha is the base for the mock server.
type Mocha struct {
	Config *Config
	T      TestingT
	Name   string

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
type TestingT interface {
	Helper()
	Logf(format string, a ...any)
	Errorf(format string, a ...any)
}

// Cleanable allows marking mocha instance to be closed on test cleanup.
type Cleanable interface {
	Cleanup(func())
}

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a ConfigBuilder implementation.
func New(t TestingT, config ...Configurer) *Mocha {
	if t == nil {
		t = notifier.NewConsole()
	}

	conf := defaultConfig()
	for _, configurer := range config {
		configurer.Apply(conf)
	}

	ctx, cancel := context.WithCancel(context.Background())
	store := newStore()
	events := newEvents()

	recovery := &recoverMid{d: conf.Debug, t: t, evt: events}

	parsers := make([]RequestBodyParser, 0, len(conf.RequestBodyParsers)+4)
	parsers = append(parsers, conf.RequestBodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &bytesParser{})

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, recovery.Recover)

	if conf.LogLevel > LogSilently {
		h := newInternalEvents(t, conf.LogLevel)

		_ = events.Subscribe(EventOnRequest, h.OnRequest)
		_ = events.Subscribe(EventOnRequestMatched, h.OnRequestMatched)
		_ = events.Subscribe(EventOnRequestNotMatched, h.OnRequestNotMatched)
		_ = events.Subscribe(EventOnError, h.OnError)
	}

	if conf.CORS != nil {
		middlewares = append(middlewares, corsMid(conf.CORS))
	}

	middlewares = append(middlewares, conf.Middlewares...)

	params := reply.Parameters()
	if conf.Parameters != nil {
		params = conf.Parameters
	}

	var p *proxy
	if conf.Proxy != nil {
		p = newProxy(conf.Proxy, events)
	}

	handler := mid.
		Compose(middlewares...).
		Root(newHandler(conf.Name, store, parsers, params, p, events, t, conf.Debug))

	if conf.HandlerDecorator != nil {
		handler = conf.HandlerDecorator(handler)
	}

	server := conf.Server

	if server == nil {
		server = newServer()
	}

	err := server.Configure(conf, handler)
	if err != nil {
		t.Logf("failed to configure server. reason=%v", err)
		panic(err)
	}

	loaders := make([]Loader, len(conf.Loaders)+1 /* number of internal loaders */)
	loaders[0] = &FileLoader{}
	for i, loader := range conf.Loaders {
		loaders[i+1] = loader
	}

	return &Mocha{
		Config: conf,
		T:      t,

		server:  server,
		storage: store,
		ctx:     ctx,
		cancel:  cancel,
		params:  params,
		scopes:  make([]*Scoped, 0),
		events:  events,
		loaders: loaders,
	}
}

// Default creates a new mock server with default configurations.
func Default() *Mocha {
	return New(notifier.NewConsole())
}

// Start starts the mock server.
func (m *Mocha) Start() (ServerInfo, error) {
	info, err := m.server.Start()
	if err != nil {
		return ServerInfo{}, err
	}

	err = m.onStart()
	if err != nil {
		return ServerInfo{}, err
	}

	return info, nil
}

// MustStart starts the mock server.
// It fails immediately if any error occurs.
func (m *Mocha) MustStart() ServerInfo {
	info, err := m.Start()
	if err != nil {
		m.T.Logf("failed to start mock server. reason=%v", err)
		panic(err)
	}

	return info
}

// StartTLS starts TLS on a mock server.
func (m *Mocha) StartTLS() (ServerInfo, error) {
	info, err := m.server.StartTLS()
	if err != nil {
		return ServerInfo{}, err
	}

	err = m.onStart()
	if err != nil {
		return ServerInfo{}, err
	}

	return info, nil
}

// MustStartTLS starts TLS on a mock server.
// It fails immediately if any error occurs.
func (m *Mocha) MustStartTLS() ServerInfo {
	info, err := m.server.StartTLS()
	if err != nil {
		m.T.Logf("failed to start a TLS mock server. reason=%v", err)
		panic(err)
	}

	return info
}

// Mock adds one or multiple request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checking if they were called or not.
//
// Usage:
//
//	scoped := m.MustMock(
//		Get(matcher.URLPath("/test")).
//			Header("test", matcher.Equal("hello")).
//			Query("filter", matcher.Equal("all")).
//			Reply(reply.Created().PlainText("hello world")))
//
//	assert.True(T, scoped.Called())
func (m *Mocha) Mock(builders ...Builder) (*Scoped, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := len(builders)
	added := make([]*Mock, size)

	for i, b := range builders {
		nm, err := b.Build()
		if err != nil {
			return nil, fmt.Errorf("error building mock at index [%d]. %w", i, err)
		}

		m.storage.Save(nm)
		added[i] = nm
	}

	scoped := scope(m.storage, added)
	m.scopes = append(m.scopes, scoped)

	return scoped, nil
}

// MustMock adds one or multiple request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checking if they were called or not.
// It fails immediately if any error occurs.
//
// Usage:
//
//	scoped := m.MustMock(
//		Get(matcher.URLPath("/test")).
//			Header("test", matcher.Equal("hello")).
//			Query("filter", matcher.Equal("all")).
//			Reply(reply.Created().PlainText("hello world")))
//
//	assert.True(T, scoped.Called())
func (m *Mocha) MustMock(builders ...Builder) *Scoped {
	scoped, err := m.Mock(builders...)
	if err != nil {
		m.T.Logf(err.Error())
		panic(err)
	}

	return scoped
}

// Parameters returns an editable parameters reply.Params that will be available when build a reply.Reply.
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
func (m *Mocha) Subscribe(evt reflect.Type, fn func(payload any)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.events.Subscribe(evt, fn)
}

// MustSubscribe add a new event listener.
// Panics if any errors occur.
func (m *Mocha) MustSubscribe(evt reflect.Type, fn func(payload any)) {
	err := m.Subscribe(evt, fn)
	if err != nil {
		panic(err)
	}
}

// Reload reloads mock set external sources, like Loader.
// Coded mocks will be kept.
func (m *Mocha) Reload() error {
	// remove mocks set by Loaders and then, reload, keeping the ones set via code.
	m.storage.DeleteExternal()

	for _, loader := range m.loaders {
		err := loader.Load(m)
		if err != nil {
			return err
		}
	}

	return nil
}

// MustReload reloads mock set external sources, like Loader.
// Coded mocks will be kept.
// It panics if any error occurs.
func (m *Mocha) MustReload() {
	err := m.Reload()

	if err != nil {
		m.T.Logf("error rebuild mock definitions. reason=%v", err.Error())
		panic(err)
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

// CloseWithT register Server Close function on TestingT Cleanup().
// Useful to close the server when tests finishes.
func (m *Mocha) CloseWithT(t Cleanable) *Mocha {
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

// Clean removes all scoped mocks.
func (m *Mocha) Clean() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.scopes {
		s.Clean()
	}
}

// --
// Mock Builders (Syntax Sugar)
// --

func (m *Mocha) GET(matcher matcher.Matcher) *MockBuilder {
	return Request().URL(matcher).Method(http.MethodGet)
}

// --
// Assertions
// --

// AssertCalled asserts that all mocks associated with this instance were called at least once.
func (m *Mocha) AssertCalled(t TestingT) bool {
	t.Helper()

	result := true

	for i, s := range m.scopes {
		if s.IsPending() {
			t.Logf("\nscope [%d]\n", i)
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
			t.Logf("\nscope [%d]\n", i)
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

// --
// Internals
// --

func (m *Mocha) onStart() error {
	err := m.Reload()
	if err != nil {
		return err
	}

	m.events.StartListening(m.ctx)

	return nil
}
