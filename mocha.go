package mocha

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mid"
	"github.com/vitorsalgado/mocha/v3/internal/mid/recover"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/matcher"
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
	name               string
	config             *Config
	server             Server
	storage            mockStore
	ctx                context.Context
	cancel             context.CancelFunc
	requestBodyParsers []RequestBodyParser
	params             Params
	scopes             []*Scoped
	loaders            []Loader
	rec                *recorder
	rmu                sync.RWMutex
	proxy              *reverseProxy
	te                 TemplateEngine
	extensions         map[string]Extension
	data               map[string]any
	cz                 *colorize.Colorize
	log                *zerolog.Logger
}

// TestingT is based on testing.T and is used for assertions.
// See Assert* methods on the application instance.
type TestingT interface {
	Helper()
	Logf(format string, a ...any)
	Errorf(format string, a ...any)
	Cleanup(func())
}

// Headers
const (
	HeaderAccept              = header.Accept
	HeaderAcceptEncoding      = header.AcceptEncoding
	HeaderContentType         = header.ContentType
	HeaderContentEncoding     = header.ContentEncoding
	HeaderAllow               = header.Allow
	HeaderAuthorization       = header.Authorization
	HeaderContentDisposition  = header.ContentDisposition
	HeaderVary                = header.Vary
	HeaderOrigin              = header.Origin
	HeaderContentLength       = header.ContentLength
	HeaderConnection          = header.Connection
	HeaderTrailer             = header.Trailer
	HeaderLocation            = header.Location
	HeaderCacheControl        = header.CacheControl
	HeaderCookie              = header.Cookie
	HeaderSetCookie           = header.SetCookie
	HeaderIfModifiedSince     = header.IfModifiedSince
	HeaderLastModified        = header.LastModified
	HeaderRetryAfter          = header.RetryAfter
	HeaderUpgrade             = header.Upgrade
	HeaderWWWAuthenticate     = header.WWWAuthenticate
	HeaderServer              = header.Server
	HeaderXForwardedFor       = header.XForwardedFor
	HeaderXForwardedProto     = header.XForwardedProto
	HeaderXForwardedProtocol  = header.XForwardedProtocol
	HeaderXForwardedSsl       = header.XForwardedSsl
	HeaderXUrlScheme          = header.XUrlScheme
	HeaderXHTTPMethodOverride = header.XHTTPMethodOverride
	HeaderXRealIP             = header.XRealIP
	HeaderXRequestID          = header.XRequestID
	HeaderXCorrelationID      = header.XCorrelationID
	HeaderXRequestedWith      = header.XRequestedWith

	HeaderAccessControlRequestMethod    = header.AccessControlRequestMethod
	HeaderAccessControlAllowOrigin      = header.AccessControlAllowOrigin
	HeaderAccessControlAllowMethods     = header.AccessControlAllowMethods
	HeaderAccessControlAllowHeaders     = header.AccessControlAllowHeaders
	HeaderAccessControlExposeHeaders    = header.AccessControlExposeHeaders
	HeaderAccessControlMaxAge           = header.AccessControlMaxAge
	HeaderAccessControlAllowCredentials = header.AccessControlAllowCredentials
	HeaderAccessControlRequestHeaders   = header.AccessControlRequestHeaders

	HeaderStrictTransportSecurity         = header.StrictTransportSecurity
	HeaderXContentTypeOptions             = header.XContentTypeOptions
	HeaderXXSSProtection                  = header.XXSSProtection
	HeaderXFrameOptions                   = header.XFrameOptions
	HeaderContentSecurityPolicy           = header.ContentSecurityPolicy
	HeaderContentSecurityPolicyReportOnly = header.ContentSecurityPolicyReportOnly
	HeaderXCSRFToken                      = header.XCSRFToken
	HeaderReferrerPolicy                  = header.ReferrerPolicy
)

// MIME Types
const (
	MIMEApplicationJSON                  = mimetype.JSON
	MIMEApplicationJSONCharsetUTF8       = mimetype.JSONCharsetUTF8
	MIMETextPlain                        = mimetype.TextPlain
	MIMETextPlainCharsetUTF8             = mimetype.TextPlainCharsetUTF8
	MIMETextHTML                         = mimetype.TextHTML
	MIMETextHTMLCharsetUTF8              = mimetype.TextHTMLCharsetUTF8
	MIMETextXML                          = mimetype.TextXML
	MIMETextXMLCharsetUTF8               = mimetype.TextXMLCharsetUTF8
	MIMEFormURLEncoded                   = mimetype.FormURLEncoded
	MIMEFormURLEncodedCharsetUTF8        = mimetype.FormURLEncodedCharsetUTF8
	MIMEApplicationJavaScript            = mimetype.ApplicationJavaScript
	MIMEApplicationJavaScriptCharsetUTF8 = mimetype.ApplicationJavaScriptCharsetUTF8
	MIMEApplicationXML                   = mimetype.ApplicationXML
	MIMEApplicationXMLCharsetUTF8        = mimetype.ApplicationXMLCharsetUTF8
	MIMEApplicationProtobuf              = mimetype.ApplicationProtobuf
	MIMEMultipartForm                    = mimetype.MultipartForm
	MIMEOctetStream                      = mimetype.OctetStream
)

// New creates a new Mocha mock server with the given configurations.
// Parameter config accepts a Config or a ConfigBuilder implementation.
// If no configuration is provided, a default one will be used.
// If no port is set, it will start the server on localhost using a random port.
func New(config ...Configurer) *Mocha {
	app := &Mocha{}

	conf := defaultConfig()
	for i, configurer := range config {
		err := configurer.Apply(conf)
		if err != nil {
			panic(fmt.Errorf(
				"server: error applying configuration at index %d. config type=%s, reason=%w",
				i,
				reflect.TypeOf(configurer),
				err,
			))
		}
	}

	setLog(conf, app)

	ctx, cancel := context.WithCancel(context.Background())
	cz := &colorize.Colorize{Enabled: conf.UseDescriptiveLogger}
	store := newStore()

	parsers := make([]RequestBodyParser, 0, len(conf.RequestBodyParsers)+4)
	parsers = append(parsers, conf.RequestBodyParsers...)
	parsers = append(parsers, &jsonBodyParser{}, &plainTextParser{}, &formURLEncodedParser{}, &noopParser{})

	recovery := recover.New(func(err error) { app.log.Error().Err(err).Msg(err.Error()) },
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
		rec = newRecorder(conf.Record)
	}

	var p *reverseProxy
	if conf.Proxy != nil {
		p = newProxy(app, conf.Proxy)
	}

	var lifecycle mockHTTPLifecycle
	if conf.UseDescriptiveLogger {
		lifecycle = &builtInDescriptiveMockHTTPLifecycle{app, cz}
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

	err := server.Setup(conf, handler)
	if err != nil {
		panic(fmt.Errorf("server: setup failed. %w", err))
	}

	loaders := make([]Loader, len(conf.Loaders)+1 /* number of internal loaders */)
	loaders[0] = &fileLoader{}
	for i, loader := range conf.Loaders {
		loaders[i+1] = loader
	}

	if conf.TemplateEngine == nil {
		tmpl := newGoTemplate()
		if len(conf.TemplateFunctions) > 0 {
			tmpl.FuncMap(conf.TemplateFunctions)
		}

		app.te = tmpl
	}

	app.config = conf
	app.name = conf.Name
	app.server = server
	app.storage = store
	app.ctx = ctx
	app.cancel = cancel
	app.params = params
	app.scopes = make([]*Scoped, 0)
	app.loaders = loaders
	app.rec = rec
	app.proxy = p
	app.requestBodyParsers = parsers
	app.extensions = make(map[string]Extension)
	app.cz = cz

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
	return app.name
}

// Start starts the mock server.
func (app *Mocha) Start() (ServerInfo, error) {
	info, err := app.server.Start()
	if err != nil {
		return ServerInfo{}, err
	}

	err = app.onStart()
	if err != nil {
		return ServerInfo{}, err
	}

	return info, nil
}

// MustStart starts the mock server.
// It fails immediately if any error occurs.
func (app *Mocha) MustStart() ServerInfo {
	info, err := app.Start()
	if err != nil {
		panic(fmt.Errorf("server: start failed. %w", err))
	}

	return info
}

// StartTLS starts TLS on a mock server.
func (app *Mocha) StartTLS() (ServerInfo, error) {
	info, err := app.server.StartTLS()
	if err != nil {
		return ServerInfo{}, err
	}

	err = app.onStart()
	if err != nil {
		return ServerInfo{}, err
	}

	return info, nil
}

// MustStartTLS starts TLS on a mock server.
// It fails immediately if any error occurs.
func (app *Mocha) MustStartTLS() ServerInfo {
	info, err := app.StartTLS()
	if err != nil {
		panic(fmt.Errorf("server: failed to start server with TLS. %w", err))
	}

	return info
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
	app.rmu.Lock()
	defer app.rmu.Unlock()

	size := len(builders)
	added := make([]*Mock, size)

	for i, b := range builders {
		mock, err := b.Build(app)
		if err != nil {
			return nil, fmt.Errorf("server: error adding mock at index %d.\n%w", i, err)
		}

		mock.prepare()

		app.storage.Save(mock)
		added[i] = mock
	}

	scoped := scope(app.storage, added)
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
func (app *Mocha) URL() string {
	return app.server.Info().URL
}

// Context returns the server internal context.Context.
func (app *Mocha) Context() context.Context {
	return app.ctx
}

// Config returns a copy of the mock server Config.
func (app *Mocha) Config() *Config {
	return app.config
}

// TemplateEngine returns the TemplateEngine associated with this instance.
func (app *Mocha) TemplateEngine() TemplateEngine {
	return app.te
}

// Logger returns the zerolog.Logger of this application.
func (app *Mocha) Logger() *zerolog.Logger {
	return app.log
}

// Reload reloads mocks from external sources, like the filesystem.
// Coded mocks will be kept.
func (app *Mocha) Reload() error {
	app.storage.DeleteExternal()

	for _, loader := range app.loaders {
		err := loader.Load(app)
		if err != nil {
			return err
		}
	}

	return nil
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
		app.log.Debug().Err(err).Msg("error closing mock server")
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
	app.rmu.RLock()
	defer app.rmu.RUnlock()

	hits := 0

	for _, s := range app.scopes {
		hits += s.Hits()
	}

	return hits
}

// Enable enables all mocks.
func (app *Mocha) Enable() {
	app.rmu.Lock()
	defer app.rmu.Unlock()

	for _, scoped := range app.scopes {
		scoped.Enable()
	}
}

// Disable disables all mocks.
func (app *Mocha) Disable() {
	app.rmu.Lock()
	defer app.rmu.Unlock()

	for _, scoped := range app.scopes {
		scoped.Disable()
	}
}

// Clean removes all scoped mocks.
func (app *Mocha) Clean() {
	app.rmu.Lock()
	defer app.rmu.Unlock()

	for _, s := range app.scopes {
		s.Clean()
	}
}

func (app *Mocha) StopRecording() {
	app.rec.stop()
}

func (app *Mocha) RegisterExtension(extension Extension) error {
	app.rmu.Lock()
	defer app.rmu.Unlock()

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
	s.WriteString(strings.Join(app.config.Directories, ", "))
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

func (app *Mocha) GET(matcher matcher.Matcher) *MockBuilder {
	return Request().URL(matcher).Method(http.MethodGet)
}

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
			app.name,
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
			app.name,
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

	t.Errorf("\nServer: %s\n Expected %d matched request hits.\n Got %d", app.name, expected, hits)

	return false
}

// --
// Internals
// --

func (app *Mocha) onStart() error {
	err := app.Reload()
	if err != nil {
		return err
	}

	if app.rec != nil {
		app.rec.start(app.ctx)
	}

	return nil
}

func setLog(conf *Config, app *Mocha) {
	if conf.Logger != nil {
		app.log = conf.Logger
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
	app.log = &log
}
