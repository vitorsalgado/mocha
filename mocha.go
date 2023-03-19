package mocha

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/internal/mid"
	"github.com/vitorsalgado/mocha/v3/internal/mid/recover"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/matcher/mfeat"
)

const (
	Version = "3.0.0"
)

// StatusNoMatch describes an HTTP response where no Mock was found.
//
// It uses http.StatusTeapot to reduce the chance of using the same
// expected response from the actual server being mocked.
// Basically, every request that doesn't match against a Mock, will have a response with http.StatusTeapot.
const StatusNoMatch = http.StatusTeapot

// Mocha is the base for the mock server.
type Mocha struct {
	config             *Config
	server             Server
	storage            mockStore
	scenarioStore      *mfeat.ScenarioStore
	ctx                context.Context
	cancel             context.CancelFunc
	requestBodyParsers []RequestBodyParser
	params             Params
	scopes             []*Scoped
	loaders            []Loader
	rec                *recorder
	rwMutex            sync.RWMutex
	proxy              *reverseProxy
	templateEngine     TemplateEngine
	extensions         map[string]Extension
	data               map[string]any
	colorizer          *colorize.Colorize
	logger             *zerolog.Logger
	startOnce          sync.Once
	random             *rand.Rand
}

// TestingT is based on testing.T and is used for assertions.
// See Assert* methods on the application instance.
type TestingT interface {
	Helper()
	Logf(format string, a ...any)
	Errorf(format string, a ...any)
	Cleanup(func())
}

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a ConfigBuilder implementation.
// If no configuration is provided, a default one will be used.
// If no port is set, it will start the server on localhost using a random port.
func New(config ...Configurer) *Mocha {
	app := &Mocha{}
	app.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	ctx, cancel := context.WithCancel(context.Background())
	app.ctx = ctx
	app.cancel = cancel

	conf := defaultConfig()
	app.config = conf

	for i, configurer := range config {
		err := configurer.Apply(conf)
		if err != nil {
			panic(fmt.Errorf(
				"server: error applying configuration at index %d with type %v\n%w",
				i,
				reflect.TypeOf(configurer),
				err,
			))
		}
	}

	app.setLog(conf)

	colors := &colorize.Colorize{Enabled: conf.UseDescriptiveLogger}
	store := newStore()

	parsers := make([]RequestBodyParser, 0, len(conf.RequestBodyParsers)+3)
	parsers = append(parsers, conf.RequestBodyParsers...)
	parsers = append(parsers, &plainTextParser{}, &formURLEncodedParser{}, &noopParser{})

	recovery := recover.New(func(err error) { app.logger.Error().Err(err).Msg(err.Error()) },
		conf.RequestWasNotMatchedStatusCode)

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, recovery.Recover)

	if conf.CORS != nil {
		middlewares = append(middlewares, corsMid(conf.CORS))
	}

	middlewares = append(middlewares, conf.Middlewares...)

	var params Params
	if conf.Parameters != nil {
		params = conf.Parameters
	} else {
		params = newInMemoryParameters()
	}

	var rec *recorder
	if conf.Record != nil {
		rec = newRecorder(app, conf.Record)
	}

	var p *reverseProxy
	if conf.Proxy != nil {
		p = newProxy(app.logger, conf)
	}

	var lifecycle mockHTTPLifecycle
	if conf.UseDescriptiveLogger {
		lifecycle = &builtInDescriptiveMockHTTPLifecycle{app, colors}
	} else {
		lifecycle = &builtInMockHTTPLifecycle{app}
	}

	handler := mid.
		Compose(middlewares...).
		Root(&mockHandler{app, lifecycle})

	if conf.HandlerDecorator != nil {
		handler = conf.HandlerDecorator(handler)
	}

	server := conf.Server

	if server == nil {
		server = newServer()
	}

	loaders := make([]Loader, 0, len(conf.Loaders)+1 /* 1 is the number of internal loaders */)
	loaders = append(loaders, &fileLoader{})
	loaders = append(loaders, conf.Loaders...)

	if conf.TemplateEngine == nil {
		tmpl := newGoTemplate()
		if len(conf.TemplateFunctions) > 0 {
			tmpl.FuncMap(conf.TemplateFunctions)
		}

		app.templateEngine = tmpl
	}

	app.server = server
	app.storage = store
	app.scenarioStore = mfeat.NewScenarioStore()

	app.params = params
	app.scopes = make([]*Scoped, 0)
	app.loaders = loaders
	app.rec = rec
	app.proxy = p
	app.requestBodyParsers = parsers
	app.extensions = make(map[string]Extension)
	app.colorizer = colors

	if app.config.Forward != nil {
		app.MustMock(Request().
			MethodMatches(matcher.Anything()).
			Priority(10).
			Reply(From(app.config.Forward.Target).
				Headers(app.config.Forward.Headers).
				ProxyHeaders(app.config.Forward.ProxyHeaders).
				RemoveProxyHeaders(app.config.Forward.ProxyHeadersToRemove...).
				TrimPrefix(app.config.Forward.TrimPrefix).
				TrimSuffix(app.config.Forward.TrimSuffix).
				SSLVerify(app.config.Forward.SSLVerify)))
	}

	err := server.Setup(app, handler)
	if err != nil {
		panic(fmt.Errorf("server: setup failed. %w", err))
	}

	return app
}

// NewT creates a new Mocha mock server with the given configurations and
// closes the server when the provided TestingT instance finishes.
func NewT(t TestingT, config ...Configurer) *Mocha {
	app := New(config...)
	t.Cleanup(app.Close)

	return app
}

// Name returns mock server name.
func (app *Mocha) Name() string {
	return app.config.Name
}

// Start starts the mock server.
func (app *Mocha) Start() error {
	err := app.server.Start()
	if err != nil {
		return err
	}

	return app.onStart()
}

// MustStart starts the mock server.
// It fails immediately if any error occurs.
func (app *Mocha) MustStart() {
	err := app.Start()
	if err != nil {
		panic(fmt.Errorf("server: start failed. %w", err))
	}
}

// StartTLS starts TLS on a mock server.
func (app *Mocha) StartTLS() error {
	err := app.server.StartTLS()
	if err != nil {
		return err
	}

	return app.onStart()
}

// MustStartTLS starts TLS on a mock server.
// It fails immediately if any error occurs.
func (app *Mocha) MustStartTLS() {
	err := app.StartTLS()
	if err != nil {
		panic(fmt.Errorf("server: failed to start server with TLS. %w", err))
	}
}

// Mock adds one or multiple request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checks if they were called or not.
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
func (app *Mocha) Mock(builders ...Builder) (*Scoped, error) {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	added := make([]string, len(builders))

	for i, b := range builders {
		mock, err := b.Build(app)
		if err != nil {
			return nil, fmt.Errorf("server: error adding mock at index %d.\n%w", i, err)
		}

		mock.prepare()

		app.storage.Save(mock)
		added[i] = mock.ID
	}

	scoped := newScope(app.storage, added)
	app.scopes = append(app.scopes, scoped)

	return scoped, nil
}

// MustMock adds one or multiple request mocks.
// It returns a Scoped instance that allows control of the added mocks and also checks if they were called or not.
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
func (app *Mocha) MustMock(builders ...Builder) *Scoped {
	scoped, err := app.Mock(builders...)
	if err != nil {
		panic(err)
	}

	return scoped
}

// Parameters returns Params instance associated with this application.
func (app *Mocha) Parameters() Params {
	return app.params
}

// URL returns the base URL string for the mock server.
// If parameter paths is provided, it will be concatenated with the server base URL.
func (app *Mocha) URL(paths ...string) string {
	if len(paths) == 0 {
		return app.server.Info().URL
	}

	u, err := url.JoinPath(app.server.Info().URL, paths...)
	if err != nil {
		app.logger.Fatal().Err(err).
			Msgf("server: building server url with path elements %s", paths)
	}

	return u
}

// Context returns the server internal context.Context.
func (app *Mocha) Context() context.Context {
	return app.ctx
}

// Config returns a copy of the mock server Config.
func (app *Mocha) Config() *Config {
	return app.config
}

// Server returns the Server implementation being used by this application.
func (app *Mocha) Server() Server {
	return app.server
}

// TemplateEngine returns the TemplateEngine associated with this instance.
func (app *Mocha) TemplateEngine() TemplateEngine {
	return app.templateEngine
}

// Logger returns the zerolog.Logger of this application.
func (app *Mocha) Logger() *zerolog.Logger {
	return app.logger
}

// Reload reloads mocks from external sources, like the filesystem.
// Coded mocks will be kept.
func (app *Mocha) Reload() error {
	app.storage.DeleteExternal()
	return app.load()
}

// MustReload reloads mocks from external sources, like the filesystem.
// Coded mocks will be kept.
// It panics if any error occurs.
func (app *Mocha) MustReload() {
	err := app.Reload()

	if err != nil {
		panic(fmt.Errorf("server: error reloading mock definitions. %w", err))
	}
}

// Close closes the mock server.
func (app *Mocha) Close() {
	app.cancel()

	err := app.server.Close()
	if err != nil {
		app.logger.Debug().Err(err).Msg("error closing mock server")
	}

	if app.rec != nil {
		app.rec.stop()
	}
}

// CloseWithT register Server Close function on TestingT Cleanup().
// Useful to close the server when tests finish.
func (app *Mocha) CloseWithT(t TestingT) *Mocha {
	t.Cleanup(func() { app.Close() })
	return app
}

// Hits returns the total matched request hits.
func (app *Mocha) Hits() int {
	app.rwMutex.RLock()
	defer app.rwMutex.RUnlock()

	hits := 0

	for _, s := range app.scopes {
		hits += s.Hits()
	}

	return hits
}

// Enable enables all mocks.
func (app *Mocha) Enable() {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	for _, scoped := range app.scopes {
		scoped.Enable()
	}
}

// Disable disables all mocks.
func (app *Mocha) Disable() {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	for _, scoped := range app.scopes {
		scoped.Disable()
	}
}

// Clean removes all scoped mocks.
func (app *Mocha) Clean() {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	for _, s := range app.scopes {
		s.Clean()
	}
}

func (app *Mocha) StopRecording() {
	app.rec.stop()
}

func (app *Mocha) RegisterExtension(extension Extension) error {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	_, ok := app.extensions[extension.UniqueName()]
	if ok {
		return fmt.Errorf(
			"server: there is already an extension registered with the name \"%s\"",
			extension.UniqueName(),
		)
	}

	app.extensions[extension.UniqueName()] = extension

	return nil
}

// SetData sets the data to be used as template data during mock configurations parsing.
func (app *Mocha) SetData(data map[string]any) {
	app.data = data
}

// Data returns the template data associated to this instance.
func (app *Mocha) Data() map[string]any {
	return app.data
}

// PrintConfig prints key configurations using the given io.Writer.
func (app *Mocha) PrintConfig(w io.Writer) error {
	s := strings.Builder{}

	if app.Name() != "" {
		s.WriteString("Server Name: ")
		s.WriteString(app.Name())
		s.WriteString("\n")
	}

	s.WriteString("Mock Search Patterns: ")
	s.WriteString(strings.Join(app.config.MockFileSearchPatterns, ", "))
	s.WriteString("\n")

	s.WriteString("Log: ")
	s.WriteString(app.config.LogVerbosity.String())
	s.WriteString("\n")

	if app.config.Proxy != nil {
		s.WriteString("Reverse Proxy: ")
		s.WriteString("enabled")
		s.WriteString("\n")
	}

	if app.config.Record != nil {
		s.WriteString("Recording: ")
		s.WriteString(app.config.Record.SaveDir)
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString("Listening: ")
	s.WriteString(app.URL())
	s.WriteString("\n")

	_, err := fmt.Fprint(w, s.String())

	return err
}

// --
// Mock Builders
// --

// AnyMethod creates a new empty Builder.
func (app *Mocha) AnyMethod() *MockBuilder {
	b := &MockBuilder{mock: newMock()}
	return b.MethodMatches(matcher.Anything())
}

// Get initializes a mock for GET method.
func (app *Mocha) Get(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodGet)
}

// Getf initializes a mock for GET method.
func (app *Mocha) Getf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodGet)
}

// Post initializes a mock for Post method.
func (app *Mocha) Post(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPost)
}

// Postf initializes a mock for Post method.
func (app *Mocha) Postf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPost)
}

// Put inits a mock for Put method.
func (app *Mocha) Put(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPut)
}

// Putf initializes a mock for Put method.
func (app *Mocha) Putf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPut)
}

// Patch initializes a mock for Patch method.
func (app *Mocha) Patch(u matcher.Matcher) *MockBuilder {
	return Request().URL(u).Method(http.MethodPatch)
}

// Patchf initializes a mock for Patch method.
func (app *Mocha) Patchf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPatch)
}

// Delete initializes a mock for Delete method.
func (app *Mocha) Delete(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodDelete)
}

// Deletef initializes a mock for Delete method.
func (app *Mocha) Deletef(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodDelete)
}

// Head initializes a mock for Head method.
func (app *Mocha) Head(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodHead)
}

// Headf initializes a mock for Head method.
func (app *Mocha) Headf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodHead)
}

// --
// Assertions
// --

// AssertCalled asserts that all mocks associated with this instance were called at least once.
func (app *Mocha) AssertCalled(t TestingT) bool {
	t.Helper()

	result := true
	size := 0
	buf := strings.Builder{}

	for i, s := range app.scopes {
		if s.IsPending() {
			buf.WriteString(fmt.Sprintf("   Scope %d\n", i))
			pending := s.GetPending()
			size += len(pending)

			for _, p := range pending {
				buf.WriteString("    Mock [")
				buf.WriteString(p.ID)
				buf.WriteString("] ")
				buf.WriteString(p.Name)
				buf.WriteString("\n")
			}

			result = false
		}
	}

	if !result {
		t.Errorf("\nServer: %s\n  There are still %d mocks that were not called.\n  Pending:\n%s",
			app.Name(),
			size,
			buf.String(),
		)
	}

	return result
}

// AssertNotCalled asserts that all mocks associated with this instance were called at least once.
func (app *Mocha) AssertNotCalled(t TestingT) bool {
	t.Helper()

	result := true
	size := 0
	buf := strings.Builder{}

	for i, s := range app.scopes {
		if !s.IsPending() {
			buf.WriteString(fmt.Sprintf("   Scope %d\n", i))
			called := s.GetCalled()
			size += len(called)

			for _, p := range called {
				buf.WriteString("    Mock [")
				buf.WriteString(p.ID)
				buf.WriteString("] ")
				buf.WriteString(p.Name)
				buf.WriteString("\n")
			}

			result = false
		}
	}

	if !result {
		t.Errorf(
			"\nServer: %s\n  %d Mock(s) were called at least once when none should be.\n  Called:\n%s",
			app.Name(),
			size,
			buf.String(),
		)
	}

	return result
}

// AssertNumberOfCalls asserts that the sum of matched request hits
// is equal to the given expected value.
func (app *Mocha) AssertNumberOfCalls(t TestingT, expected int) bool {
	t.Helper()

	hits := app.Hits()

	if hits == expected {
		return true
	}

	t.Errorf("\nServer: %s\n Expected %d matched request hits.\n Got %d", app.Name(), expected, hits)

	return false
}

// --
// Internals
// --

func (app *Mocha) onStart() (err error) {
	app.startOnce.Do(func() { err = app.load() })

	if err != nil {
		return err
	}

	if app.rec != nil {
		app.rec.start(app.ctx)
	}

	return nil
}

func (app *Mocha) load() error {
	for _, loader := range app.loaders {
		err := loader.Load(app)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *Mocha) setLog(conf *Config) {
	if conf.Logger != nil {
		app.logger = conf.Logger
		return
	}

	var output io.Writer
	if conf.LogPretty {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		output = os.Stdout
	}

	c := zerolog.New(output).Level(zerolog.Level(conf.LogLevel)).With().Timestamp()
	if conf.Name != "" {
		c.Str("server", conf.Name)
	}

	log := c.Logger()
	app.logger = &log
}
