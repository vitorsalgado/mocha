package mhttp

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/internal/mid"
	"github.com/vitorsalgado/mocha/v3/internal/mid/recover"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/matcher/mfeat"
	"github.com/vitorsalgado/mocha/v3/mhttp/cors"
)

var (
	_ foundation.MockApp[*HTTPMock] = (*HTTPMockApp)(nil)
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

// HTTPMockApp is the base for the mock server.
type HTTPMockApp struct {
	*foundation.BaseApp[*HTTPMock, *HTTPMockApp]

	config             *Config
	server             Server
	storage            *foundation.MockStore[*HTTPMock]
	scenarioStore      *mfeat.ScenarioStore
	ctx                context.Context
	cancel             context.CancelFunc
	requestBodyParsers []RequestBodyParser
	params             foundation.Params
	loaders            []Loader
	rec                *recorder
	rwMutex            sync.RWMutex
	proxy              *reverseProxy
	templateEngine     foundation.TemplateEngine
	extensions         map[string]foundation.Extension
	data               map[string]any
	colorizer          *colorize.Colorize
	logger             *zerolog.Logger
	startOnce          sync.Once
	random             *rand.Rand
}

// NewAPI creates a new HTTPMockApp mock server with the given configurations.
// Parameter config accepts a Config or a ConfigBuilder implementation.
// If no configuration is provided, a default one will be used.
// If no port is set, it will start the server on localhost using a random port.
func NewAPI(config ...Configurer) *HTTPMockApp {
	app := &HTTPMockApp{}
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
				"server: error applying configuration at index %d with type %T\n%w",
				i,
				configurer,
				err,
			))
		}
	}

	app.logger = app.getLog(conf)

	colors := &colorize.Colorize{Enabled: conf.UseDescriptiveLogger}
	store := foundation.NewStore[*HTTPMock]()

	parsers := make([]RequestBodyParser, 0, len(conf.RequestBodyParsers)+3)
	parsers = append(parsers, conf.RequestBodyParsers...)
	parsers = append(parsers, &plainTextParser{}, &formURLEncodedParser{}, &noopParser{})

	recovery := recover.New(func(err error) { app.logger.Error().Err(err).Msg(err.Error()) },
		conf.RequestWasNotMatchedStatusCode)

	middlewares := make([]func(handler http.Handler) http.Handler, 0)
	middlewares = append(middlewares, recovery.Recover)

	if conf.CORS != nil {
		middlewares = append(middlewares, cors.New(conf.CORS))
	}

	middlewares = append(middlewares, conf.Middlewares...)

	var params foundation.Params
	if conf.Parameters != nil {
		params = conf.Parameters
	} else {
		params = foundation.NewInMemoryParameters()
	}

	var rec *recorder
	if conf.Record != nil {
		rec = newRecorder(app, conf.Record)
	}

	var p *reverseProxy
	if conf.Proxy != nil {
		p = newProxy(app.Logger(), conf)
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

	app.BaseApp = foundation.NewBaseApp[*HTTPMock, *HTTPMockApp](app, store)
	app.server = server
	app.storage = store
	app.scenarioStore = mfeat.NewScenarioStore()
	app.params = params
	app.loaders = loaders
	app.rec = rec
	app.proxy = p
	app.requestBodyParsers = parsers
	app.extensions = make(map[string]foundation.Extension)
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

// NewAPIWithT creates a new HTTPMockApp mock server with the given configurations and
// closes the server when the provided TestingT instance finishes.
func NewAPIWithT(t foundation.TestingT, config ...Configurer) *HTTPMockApp {
	app := NewAPI(config...)
	t.Cleanup(app.Close)

	return app
}

// Name returns mock server name.
func (app *HTTPMockApp) Name() string {
	return app.config.Name
}

// Start starts the mock server.
func (app *HTTPMockApp) Start() error {
	err := app.server.Start()
	if err != nil {
		return err
	}

	return app.onStart()
}

// MustStart starts the mock server.
// It fails immediately if any error occurs.
func (app *HTTPMockApp) MustStart() {
	err := app.Start()
	if err != nil {
		panic(fmt.Errorf("server: start failed. %w", err))
	}
}

// StartTLS starts TLS on a mock server.
func (app *HTTPMockApp) StartTLS() error {
	err := app.server.StartTLS()
	if err != nil {
		return err
	}

	return app.onStart()
}

// MustStartTLS starts TLS on a mock server.
// It fails immediately if any error occurs.
func (app *HTTPMockApp) MustStartTLS() {
	err := app.StartTLS()
	if err != nil {
		panic(fmt.Errorf("server: failed to start server with TLS. %w", err))
	}
}

// URL returns the base URL string for the mock server.
// If parameter paths is provided, it will be concatenated with the server base URL.
func (app *HTTPMockApp) URL(paths ...string) string {
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
func (app *HTTPMockApp) Context() context.Context {
	return app.ctx
}

// Parameters returns Params instance associated with this application.
func (app *HTTPMockApp) Parameters() foundation.Params {
	return app.params
}

// Config returns a copy of the mock server Config.
func (app *HTTPMockApp) Config() *Config {
	return app.config
}

// Server returns the Server implementation being used by this application.
func (app *HTTPMockApp) Server() Server {
	return app.server
}

// TemplateEngine returns the TemplateEngine associated with this instance.
func (app *HTTPMockApp) TemplateEngine() foundation.TemplateEngine {
	return app.templateEngine
}

// Logger returns the zerolog.Logger of this application.
func (app *HTTPMockApp) Logger() *zerolog.Logger {
	return app.logger
}

// Reload reloads mocks from external sources, like the filesystem.
// Coded mocks will be kept.
func (app *HTTPMockApp) Reload() error {
	// app.storage.DeleteExternal()
	return app.load()
}

// MustReload reloads mocks from external sources, like the filesystem.
// Coded mocks will be kept.
// It panics if any error occurs.
func (app *HTTPMockApp) MustReload() {
	err := app.Reload()

	if err != nil {
		panic(fmt.Errorf("server: error reloading mock definitions. %w", err))
	}
}

// Close closes the mock server.
func (app *HTTPMockApp) Close() {
	app.cancel()

	err := app.server.Close()
	if err != nil {
		app.logger.Debug().Err(err).Msg("server: Close: error closing server")
	}

	if app.rec != nil {
		app.rec.stop()
	}
}

func (app *HTTPMockApp) CloseNow() {
	app.cancel()

	err := app.server.CloseNow()
	if err != nil {
		app.logger.Debug().Err(err).Msg("server: CloseNow: error force closing server")
	}

	if app.rec != nil {
		app.rec.stop()
	}
}

// CloseWithT register Server Close function on TestingT Cleanup().
// Useful to close the server when tests finish.
func (app *HTTPMockApp) CloseWithT(t foundation.TestingT) *HTTPMockApp {
	t.Cleanup(func() { app.Close() })
	return app
}

func (app *HTTPMockApp) StopRecording() {
	app.rec.stop()
}

func (app *HTTPMockApp) RegisterExtension(extension foundation.Extension) error {
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
func (app *HTTPMockApp) SetData(data map[string]any) {
	app.data = data
}

// Data returns the template data associated to this instance.
func (app *HTTPMockApp) Data() map[string]any {
	return app.data
}

// PrintConfig prints key configurations using the given io.Writer.
func (app *HTTPMockApp) PrintConfig(w io.Writer) error {
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
// Mock Builder Initializers
// --

// AnyMethod creates a new empty Builder.
func (app *HTTPMockApp) AnyMethod() *HTTPMockBuilder {
	b := &HTTPMockBuilder{mock: newMock()}
	return b.MethodMatches(matcher.Anything())
}

// Get initializes a mock for GET method.
func (app *HTTPMockApp) Get(m matcher.Matcher) *HTTPMockBuilder {
	return Request().URL(m).Method(http.MethodGet)
}

// Getf initializes a mock for GET method.
func (app *HTTPMockApp) Getf(path string, a ...any) *HTTPMockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodGet)
}

// Post initializes a mock for Post method.
func (app *HTTPMockApp) Post(m matcher.Matcher) *HTTPMockBuilder {
	return Request().URL(m).Method(http.MethodPost)
}

// Postf initializes a mock for Post method.
func (app *HTTPMockApp) Postf(path string, a ...any) *HTTPMockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPost)
}

// Put inits a mock for Put method.
func (app *HTTPMockApp) Put(m matcher.Matcher) *HTTPMockBuilder {
	return Request().URL(m).Method(http.MethodPut)
}

// Putf initializes a mock for Put method.
func (app *HTTPMockApp) Putf(path string, a ...any) *HTTPMockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPut)
}

// Patch initializes a mock for Patch method.
func (app *HTTPMockApp) Patch(u matcher.Matcher) *HTTPMockBuilder {
	return Request().URL(u).Method(http.MethodPatch)
}

// Patchf initializes a mock for Patch method.
func (app *HTTPMockApp) Patchf(path string, a ...any) *HTTPMockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPatch)
}

// Delete initializes a mock for Delete method.
func (app *HTTPMockApp) Delete(m matcher.Matcher) *HTTPMockBuilder {
	return Request().URL(m).Method(http.MethodDelete)
}

// Deletef initializes a mock for Delete method.
func (app *HTTPMockApp) Deletef(path string, a ...any) *HTTPMockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodDelete)
}

// Head initializes a mock for Head method.
func (app *HTTPMockApp) Head(m matcher.Matcher) *HTTPMockBuilder {
	return Request().URL(m).Method(http.MethodHead)
}

// Headf initializes a mock for Head method.
func (app *HTTPMockApp) Headf(path string, a ...any) *HTTPMockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodHead)
}

// --
// Internals
// --

func (app *HTTPMockApp) onStart() (err error) {
	app.startOnce.Do(func() { err = app.load() })

	if err != nil {
		return err
	}

	if app.rec != nil {
		app.rec.start(app.ctx)
	}

	return nil
}

func (app *HTTPMockApp) load() error {
	for _, loader := range app.loaders {
		err := loader.Load(app)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *HTTPMockApp) getLog(conf *Config) *zerolog.Logger {
	if conf.Logger != nil {
		return conf.Logger
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
	return &log
}
