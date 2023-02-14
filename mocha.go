package mocha

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/internal/logger"
	"github.com/vitorsalgado/mocha/v3/internal/mid"
	"github.com/vitorsalgado/mocha/v3/internal/mid/recover"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/x/event"
)

const (
	Version = "3.0.0"
)

// Mocha is the base for the mock server.
type Mocha struct {
	log                logger.Log
	name               string
	config             *Config
	server             Server
	storage            mockStore
	ctx                context.Context
	cancel             context.CancelFunc
	requestBodyParsers []RequestBodyParser
	params             Params
	listener           *event.Listener
	scopes             []*Scoped
	loaders            []Loader
	rec                *record
	rmu                sync.RWMutex
	rmuMock            sync.RWMutex
	proxy              *reverseProxy
	extensions         map[string]Extension
}

// TestingT is based on testing.T and allow mocha components to log information and errors.
type TestingT interface {
	Helper()
	Logf(format string, a ...any)
	Errorf(format string, a ...any)
	Cleanup(func())
}

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a ConfigBuilder implementation.
func New(config ...Configurer) *Mocha {
	m := &Mocha{}
	l := logger.NewConsole()

	conf := defaultConfig()
	for i, configurer := range config {
		err := configurer.Apply(conf)
		if err != nil {
			l.Logf("error applying configuration [%d]. reason=%s", i, err.Error())
			panic(err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	store := newStore()
	events := event.New()

	parsers := make([]RequestBodyParser, 0, len(conf.RequestBodyParsers)+4)
	parsers = append(parsers, conf.RequestBodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &noopParser{})

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, recover.New(l).Recover)

	if conf.LogLevel > LogSilently {
		h := event.NewInternalListener(l, conf.LogLevel == LogVerbose)

		_ = events.Subscribe(event.EventOnRequest, h.OnRequest)
		_ = events.Subscribe(event.EventOnRequestMatched, h.OnRequestMatched)
		_ = events.Subscribe(event.EventOnRequestNotMatched, h.OnRequestNotMatched)
		_ = events.Subscribe(event.EventOnError, h.OnError)
	}

	if conf.CORS != nil {
		middlewares = append(middlewares, corsMid(conf.CORS))
	}

	middlewares = append(middlewares, conf.Middlewares...)

	params := newInMemoryParameters()
	if conf.Parameters != nil {
		params = conf.Parameters
	}

	var rec *record
	if conf.Record != nil {
		rec = newRecorder(conf.Record)
	}

	var p *reverseProxy
	if conf.Proxy != nil {
		p = newProxy(conf.Proxy, events)
	}

	handler := mid.
		Compose(middlewares...).
		Root(&mockHandler{m})

	if conf.HandlerDecorator != nil {
		handler = conf.HandlerDecorator(handler)
	}

	server := conf.Server

	if server == nil {
		server = newServer()
	}

	err := server.Setup(conf, handler)
	if err != nil {
		l.Logf("failed to configure server. reason=%v", err)
		panic(err)
	}

	loaders := make([]Loader, len(conf.Loaders)+1 /* number of internal loaders */)
	loaders[0] = &FileLoader{}
	for i, loader := range conf.Loaders {
		loaders[i+1] = loader
	}

	m.config = conf
	m.log = l
	m.name = conf.Name
	m.server = server
	m.storage = store
	m.ctx = ctx
	m.cancel = cancel
	m.params = params
	m.listener = events
	m.scopes = make([]*Scoped, 0)
	m.loaders = loaders
	m.rec = rec
	m.proxy = p
	m.requestBodyParsers = parsers
	m.extensions = make(map[string]Extension)

	if m.config.Forward != nil {
		m.MustMock(Request().
			MethodMatches(matcher.Anything()).
			Priority(10).
			Reply(From(m.config.Forward.Target).
				Headers(m.config.Forward.Headers).
				ProxyHeaders(m.config.Forward.ProxyHeaders).
				RemoveProxyHeaders(m.config.Forward.ProxyHeadersToRemove...).
				TrimPrefix(m.config.Forward.TrimPrefix).
				TrimSuffix(m.config.Forward.TrimSuffix)))
	}

	return m
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
		m.log.Logf("failed to start mock server. reason=%v", err)
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
	info, err := m.StartTLS()
	if err != nil {
		m.log.Logf("failed to start a TLS mock server. reason=%v", err)
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
//			Header("test", matcher.StrictEqual("hello")).
//			Query("filter", matcher.StrictEqual("all")).
//			Reply(reply.Created().PlainText("hello world")))
//
//	assert.True(txtTemplate, scoped.HasBeenCalled())
func (m *Mocha) Mock(builders ...Builder) (*Scoped, error) {
	m.rmuMock.Lock()
	defer m.rmuMock.Unlock()

	size := len(builders)
	added := make([]*Mock, size)

	for i, b := range builders {
		mock, err := b.Build(m)
		if err != nil {
			return nil, fmt.Errorf("error building mock at index [%d]. %w", i, err)
		}

		mock.prepare()

		m.storage.Save(mock)
		added[i] = mock
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
//			Header("test", matcher.StrictEqual("hello")).
//			Query("filter", matcher.StrictEqual("all")).
//			Reply(reply.Created().PlainText("hello world")))
//
//	assert.True(txtTemplate, scoped.HasBeenCalled())
func (m *Mocha) MustMock(builders ...Builder) *Scoped {
	scoped, err := m.Mock(builders...)
	if err != nil {
		m.log.Logf(err.Error())
		panic(err)
	}

	return scoped
}

// Parameters returns an editable parameters reply.Params that will be available when build a reply.Reply.
func (m *Mocha) Parameters() Params {
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

// Config returns the mock server configurations.
func (m *Mocha) Config() Config {
	return *m.config
}

// Name returns mock server name.
func (m *Mocha) Name() string {
	return m.name
}

// Subscribe add a new event listener.
func (m *Mocha) Subscribe(evt reflect.Type, fn func(payload any)) error {
	return m.listener.Subscribe(evt, fn)
}

// MustSubscribe add a new event listener.
// Panics if any errors occur.
func (m *Mocha) MustSubscribe(evt reflect.Type, fn func(payload any)) {
	err := m.Subscribe(evt, fn)
	if err != nil {
		panic(err)
	}
}

// Reload reloads mocks from external sources, like Loader.
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

// MustReload reloads mocks from external sources, like Loader.
// Coded mocks will be kept.
// It panics if any error occurs.
func (m *Mocha) MustReload() {
	err := m.Reload()

	if err != nil {
		m.log.Logf("error rebuild mock definitions. reason=%v", err.Error())
		panic(err)
	}
}

// Close closes the mock server.
func (m *Mocha) Close() {
	m.cancel()

	err := m.server.Close()
	if err != nil {
		m.log.Logf(err.Error())
	}

	if m.rec != nil {
		m.rec.stop()
	}
}

// CloseWithT register Server Close function on TestingT Cleanup().
// Useful to close the server when tests finishes.
func (m *Mocha) CloseWithT(t TestingT) *Mocha {
	t.Cleanup(func() { m.Close() })
	return m
}

// Hits returns the total matched request hits.
func (m *Mocha) Hits() int {
	m.rmu.RLock()
	defer m.rmu.RUnlock()

	hits := 0

	for _, s := range m.scopes {
		hits += s.Hits()
	}

	return hits
}

// Enable enables all mocks.
func (m *Mocha) Enable() {
	m.rmu.Lock()
	defer m.rmu.Unlock()

	for _, scoped := range m.scopes {
		scoped.Enable()
	}
}

// Disable disables all mocks.
func (m *Mocha) Disable() {
	m.rmu.Lock()
	defer m.rmu.Unlock()

	for _, scoped := range m.scopes {
		scoped.Disable()
	}
}

// Clean removes all scoped mocks.
func (m *Mocha) Clean() {
	m.rmuMock.Lock()
	defer m.rmuMock.Unlock()

	for _, s := range m.scopes {
		s.Clean()
	}
}

func (m *Mocha) StopRecording() {
	m.rec.stop()
}

func (m *Mocha) RegisterExtension(extension Extension) error {
	m.rmu.Lock()
	defer m.rmu.Unlock()

	_, ok := m.extensions[extension.UniqueName()]
	if ok {
		return fmt.Errorf("there is already an extension registered with the name \"%s\"", extension.UniqueName())
	}

	m.extensions[extension.UniqueName()] = extension

	return nil
}

// PrintConfig prints key configurations using the given io.Writer.
func (m *Mocha) PrintConfig(w io.Writer) error {
	s := strings.Builder{}

	if m.Name() != "" {
		s.WriteString(colorize.Bold("Name: "))
		s.WriteString(m.Name())
		s.WriteString("\n")
	}

	s.WriteString(colorize.Bold("Mock Search Patterns: "))
	s.WriteString(strings.Join(m.config.Directories, ", "))
	s.WriteString("\n")

	s.WriteString(colorize.Bold("Log: "))
	s.WriteString(m.config.LogLevel.String())
	s.WriteString("\n")

	if m.config.Proxy != nil {
		s.WriteString(colorize.Bold("Reverse Proxy: "))
		s.WriteString("enabled")
		s.WriteString("\n")
	}

	if m.config.Record != nil {
		s.WriteString(colorize.Bold("Recording: "))
		s.WriteString(m.config.Record.SaveDir)
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(colorize.Green("Listening: "))
	s.WriteString(m.URL())
	s.WriteString("\n")

	_, err := fmt.Fprint(w, s.String())

	return err
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

// AssertNumberOfCalls asserts that the sum of matched request hits
// is equal to the given expected value.
func (m *Mocha) AssertNumberOfCalls(t TestingT, expected int) bool {
	t.Helper()

	hits := m.Hits()

	if hits == expected {
		return true
	}

	t.Errorf("\nexpected [%d] matched request hits. got [%d]", expected, hits)

	return false
}

// --
// Internals
// --

func (m *Mocha) onStart() error {
	err := m.Reload()
	if err != nil {
		return err
	}

	m.listener.StartListening(m.ctx)

	if m.rec != nil {
		m.rec.start(m.ctx)
	}

	return nil
}
