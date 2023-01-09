package mocha

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
)

type LogLevel int

const (
	// LogSilently enable minimum log mode.
	LogSilently LogLevel = iota
	// LogInfo logs only informative messages, without too much Details.
	LogInfo
	// LogVerbose logs detailed information about requests, matches and non-matches.
	LogVerbose
)

func (l LogLevel) String() string {
	switch l {
	case LogSilently:
		return "silent"
	case LogInfo:
		return "info"
	default:
		return "verbose"
	}
}

// Defaults
const (
	// ConfigMockFilePattern is the default filename glob pattern to search for local mock files.
	ConfigMockFilePattern = "testdata/*mock.json"

	// StatusRequestWasNotMatch describes an HTTP response where no Mock was found.
	//
	// It uses http.StatusTeapot to reduce the chance of using the same
	// expected response from the actual server being mocked.
	// Basically, every request that doesn't match against to a Mock will return http.StatusTeapot.
	StatusRequestWasNotMatch = http.StatusTeapot
)

// Configurer lets users configure the Mock API.
type Configurer interface {
	// Apply applies a configuration.
	Apply(conf *Config) error
}

// Config holds Mocha mock server configurations.
type Config struct {
	// Name sets a name to the mock server.
	// Adds more context for when you have more mocks APIs configured.
	Name string

	// Addr defines a custom server address.
	Addr string

	// MockNotFoundStatusCode defines the status code that should be used when
	// an HTTP request doesn't match with any mock.
	// Defaults to 418 (I'm a teapot).
	MockNotFoundStatusCode int

	// RequestBodyParsers defines request body parsers to be executed before core parsers.
	RequestBodyParsers []RequestBodyParser

	// Middlewares defines a list of custom middlewares that will be
	// set after panic recover and before mock handler.
	Middlewares []func(http.Handler) http.Handler

	// CORS defines CORS configurations.
	CORS *CORSConfig

	// Server defines a custom mock HTTP server.
	Server Server

	// HandlerDecorator provides a mean to configure a custom HTTP handler
	// while leveraging the default mock handler.
	HandlerDecorator func(handler http.Handler) http.Handler

	// LogLevel defines the level of logs
	LogLevel LogLevel

	// Parameters sets a custom reply parameters store.
	Parameters Params

	// Directories configures glob patterns to load mock from the file system.
	Directories []string

	// Loaders configures additional loaders.
	Loaders []Loader

	// Proxy configures the mock server as a proxy.
	Proxy *ProxyConfig

	// Record configures Mock Request/Stub recording.
	// Needs to be used with Proxy.
	Record *RecordConfig

	// CLI Only Options

	// UseHTTPS defines that the mock server should use HTTPS.
	// This is only used running the command-line version.
	// To start an HTTPS server from code, call .StartTLS() or .MustStartTLS() from Moai instance.
	UseHTTPS bool

	// Forward configures a forward proxy for matched requests.
	Forward *ForwardConfig
}

// ForwardConfig configures a forward proxy for matched requests.
// Only for CLI.
type ForwardConfig struct {
	Target               string
	Headers              http.Header
	ProxyHeaders         http.Header
	ProxyHeadersToRemove []string
	TrimPrefix           string
	TrimSuffix           string
}

// Apply copies the current Config struct values to the given Config parameter.
// It allows the Config struct to be used as a Configurer.
func (c *Config) Apply(conf *Config) error {
	conf.Name = c.Name
	conf.Addr = c.Addr
	conf.MockNotFoundStatusCode = c.MockNotFoundStatusCode
	conf.RequestBodyParsers = c.RequestBodyParsers
	conf.Middlewares = c.Middlewares
	conf.CORS = c.CORS
	conf.Server = c.Server
	conf.HandlerDecorator = c.HandlerDecorator
	conf.LogLevel = c.LogLevel
	conf.Parameters = c.Parameters
	conf.Directories = c.Directories
	conf.Loaders = c.Loaders
	conf.Proxy = c.Proxy
	conf.Record = c.Record
	conf.Forward = c.Forward
	conf.UseHTTPS = c.UseHTTPS

	return nil
}

// configFunc is a helper to build Configurer instances with functions.
type configFunc func(config *Config)

func (f configFunc) Apply(config *Config) error {
	f(config)
	return nil
}

// ConfigBuilder lets users create a Config with a fluent API.
type ConfigBuilder struct {
	conf *Config
}

func defaultConfig() *Config {
	return &Config{
		MockNotFoundStatusCode: StatusRequestWasNotMatch,
		LogLevel:               LogVerbose,
		Directories:            []string{ConfigMockFilePattern},
		RequestBodyParsers:     make([]RequestBodyParser, 0),
		Middlewares:            make([]func(http.Handler) http.Handler, 0),
		Loaders:                make([]Loader, 0),
	}
}

// Configure initialize a new ConfigBuilder.
// Entrypoint to start a new custom configuration for Mocha mock servers.
func Configure() *ConfigBuilder {
	return &ConfigBuilder{conf: defaultConfig()}
}

// Name sets a name to the mock server.
func (cb *ConfigBuilder) Name(name string) *ConfigBuilder {
	cb.conf.Name = name
	return cb
}

// Addr sets a custom address for the mock HTTP server.
func (cb *ConfigBuilder) Addr(addr string) *ConfigBuilder {
	cb.conf.Addr = addr
	return cb
}

// MockNotFoundStatusCode defines the status code to be used no mock matches with an HTTP request.
func (cb *ConfigBuilder) MockNotFoundStatusCode(code int) *ConfigBuilder {
	cb.conf.MockNotFoundStatusCode = code
	return cb
}

// RequestBodyParsers adds a custom list of RequestBodyParsers.
func (cb *ConfigBuilder) RequestBodyParsers(bp ...RequestBodyParser) *ConfigBuilder {
	cb.conf.RequestBodyParsers = append(cb.conf.RequestBodyParsers, bp...)
	return cb
}

// Middlewares adds custom middlewares to the mock server.
// Use this to add custom request interceptors.
func (cb *ConfigBuilder) Middlewares(fn ...func(handler http.Handler) http.Handler) *ConfigBuilder {
	cb.conf.Middlewares = append(cb.conf.Middlewares, fn...)
	return cb
}

// CORS configures Cross Origin Resource Sharing for the mock server.
func (cb *ConfigBuilder) CORS(options ...CORSConfigurer) *ConfigBuilder {
	opts := &_defaultCORSConfig
	for _, option := range options {
		option.Apply(opts)
	}

	cb.conf.CORS = opts

	return cb
}

// Server configures a custom HTTP mock Server.
func (cb *ConfigBuilder) Server(srv Server) *ConfigBuilder {
	cb.conf.Server = srv
	return cb
}

// HandlerDecorator configures a custom HTTP handler using the default mock handler.
func (cb *ConfigBuilder) HandlerDecorator(fn func(handler http.Handler) http.Handler) *ConfigBuilder {
	cb.conf.HandlerDecorator = fn
	return cb
}

// LogLevel configure the verbosity of informative logs.
// Defaults to LogVerbose.
func (cb *ConfigBuilder) LogLevel(l LogLevel) *ConfigBuilder {
	cb.conf.LogLevel = l
	return cb
}

// Parameters sets a custom reply parameters store.
func (cb *ConfigBuilder) Parameters(params Params) *ConfigBuilder {
	cb.conf.Parameters = params
	return cb
}

// Dirs sets a custom Glob patterns to load mock from the file system.
// Defaults to [testdata/*.mock.json, testdata/*.mock.yaml].
func (cb *ConfigBuilder) Dirs(patterns ...string) *ConfigBuilder {
	dirs := make([]string, 0)
	dirs = append(dirs, ConfigMockFilePattern)
	dirs = append(dirs, patterns...)

	cb.conf.Directories = dirs
	return cb
}

// NewDirs configures directories to search for local mocks,
// overwriting the default internal mock filename pattern.
// Pass a list of glob patterns supported by Go Standard Library.
// Use Dirs to keep the default internal pattern.
func (cb *ConfigBuilder) NewDirs(patterns ...string) *ConfigBuilder {
	cb.conf.Directories = patterns
	return cb
}

// Loader configures an additional Loader.
func (cb *ConfigBuilder) Loader(loader Loader) *ConfigBuilder {
	cb.conf.Loaders = append(cb.conf.Loaders, loader)
	return cb
}

// Proxy configures the mock server as a proxy server.
// Non-Matched requests will be routed based on the proxy configuration.
func (cb *ConfigBuilder) Proxy(options ...ProxyConfigurer) *ConfigBuilder {
	opts := &ProxyConfig{}

	for _, option := range options {
		option.Apply(opts)
	}

	cb.conf.Proxy = opts

	return cb
}

// Record configures recording.
func (cb *ConfigBuilder) Record(options ...RecordConfigurer) *ConfigBuilder {
	opts := defaultRecordConfig()

	for _, option := range options {
		option.Apply(opts)
	}

	cb.conf.Record = opts

	return cb
}

// Apply builds a new Config with previously configured values.
func (cb *ConfigBuilder) Apply(conf *Config) error {
	return cb.conf.Apply(conf)
}

// --
// config Functions
// --

// WithName sets a name to the mock server.
func WithName(name string) Configurer {
	return configFunc(func(c *Config) { c.Name = name })
}

// WithAddr configures the server address.
func WithAddr(addr string) Configurer {
	return configFunc(func(c *Config) { c.Addr = addr })
}

// WithMockNotFoundStatusCode defines the status code to be used no mock matches with an HTTP request.
func WithMockNotFoundStatusCode(code int) Configurer {
	return configFunc(func(c *Config) { c.MockNotFoundStatusCode = code })
}

// WithRequestBodyParsers configures one or more RequestBodyParser.
func WithRequestBodyParsers(parsers ...RequestBodyParser) Configurer {
	return configFunc(func(c *Config) { c.RequestBodyParsers = append(c.RequestBodyParsers, parsers...) })
}

// WithMiddlewares adds one or more middlewares to be executed before the mock HTTP handler.
func WithMiddlewares(middlewares ...func(handler http.Handler) http.Handler) Configurer {
	return configFunc(func(c *Config) { c.Middlewares = append(c.Middlewares, middlewares...) })
}

// WithCORS configures CORS.
func WithCORS(opts ...CORSConfigurer) Configurer {
	return configFunc(func(c *Config) {
		options := &_defaultCORSConfig
		for _, option := range opts {
			option.Apply(options)
		}

		c.CORS = options
	})
}

// WithServer configures a custom mock HTTP Server.
// If none is set, a default one will be used.
func WithServer(srv Server) Configurer {
	return configFunc(func(c *Config) { c.Server = srv })
}

// WithHandlerDecorator configures a http.Handler that decorates the internal mock HTTP handler.
func WithHandlerDecorator(fn func(handler http.Handler) http.Handler) Configurer {
	return configFunc(func(c *Config) { c.HandlerDecorator = fn })
}

// WithLogLevel sets the mock server LogLevel.
func WithLogLevel(level LogLevel) Configurer {
	return configFunc(func(c *Config) { c.LogLevel = level })
}

// WithParams configures a custom reply.Params.
func WithParams(params Params) Configurer {
	return configFunc(func(c *Config) { c.Parameters = params })
}

// WithDirs configures directories to search for local mocks.
// Pass a list of glob patterns supported by Go Standard Library.
// This method keeps the default mock filename pattern, [testdata/*mock.json].
// to overwrite the default mock filename pattern, use WithNewDirs.
func WithDirs(patterns ...string) Configurer {
	return configFunc(func(c *Config) {
		dirs := make([]string, 0)
		dirs = append(dirs, ConfigMockFilePattern)
		dirs = append(dirs, patterns...)

		c.Directories = dirs
	})
}

// WithNewDirs configures directories to search for local mocks,
// overwriting the default internal mock filename pattern.
// Pass a list of glob patterns supported by Go Standard Library.
// Use WithDirs to keep the default internal pattern.
func WithNewDirs(patterns ...string) Configurer {
	return configFunc(func(c *Config) { c.Directories = patterns })
}

// WithLoader adds a new Loader to the configuration.
func WithLoader(loader Loader) Configurer {
	return configFunc(func(c *Config) { c.Loaders = append(c.Loaders, loader) })
}

// WithProxy configures the mock server as a proxy server.
func WithProxy(options ...ProxyConfigurer) Configurer {
	return configFunc(func(c *Config) {
		opts := &ProxyConfig{}

		for _, option := range options {
			option.Apply(opts)
		}

		c.Proxy = opts
	})
}

// --
// Globals
// --

// SetColors enable/disable terminal colors.
func SetColors(value bool) {
	colorize.UseColors(value)
}
