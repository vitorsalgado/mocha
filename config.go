package mocha

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/reply"
)

type LogLevel int

const (
	// LogSilently disable all logs.
	LogSilently LogLevel = iota
	// LogInfo logs only informative messages, without too much details.
	LogInfo
	// LogVerbose logs detailed information about requests, matches and non-matches.
	LogVerbose
)

// Defaults
const (
	// ConfigMockFilePattern is the default filename glob pattern to search for local mock files.
	ConfigMockFilePattern = "testdata/*mock.json"
)

// Configurer lets users configure the Mock API.
type Configurer interface {
	// Apply applies a configuration.
	Apply(conf *Config)
}

// Debug is a function to help debug server unexpected errors.
type Debug func(err error)

// Config holds Mocha mock server configurations.
type Config struct {
	// Addr defines a custom server address.
	Addr string

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
	Parameters reply.Params

	// Files configures glob patterns to load mock from the file system.
	Files []string

	Debug Debug
}

// Apply copies the current Config struct values to the given Config parameter.
// It allows the Config struct to be used as a Configurer.
func (c *Config) Apply(conf *Config) {
	conf.Addr = c.Addr
	conf.RequestBodyParsers = c.RequestBodyParsers
	conf.Middlewares = c.Middlewares
	conf.CORS = c.CORS
	conf.Server = c.Server
	conf.HandlerDecorator = c.HandlerDecorator
	conf.LogLevel = c.LogLevel
	conf.Parameters = c.Parameters
	conf.Files = c.Files
	conf.Debug = c.Debug
}

// configFunc is a helper to build config functions.
type configFunc func(config *Config)

func (f configFunc) Apply(config *Config) { f(config) }

// ConfigBuilder lets users create a Config with a fluent API.
type ConfigBuilder struct {
	conf *Config
}

func defaultConfig() *Config {
	return &Config{
		LogLevel:           LogVerbose,
		RequestBodyParsers: make([]RequestBodyParser, 0),
		Files:              []string{ConfigMockFilePattern},
		Middlewares:        make([]func(http.Handler) http.Handler, 0)}
}

// Configure inits a new ConfigBuilder.
// Entrypoint to start a new custom configuration for Mocha mock servers.
func Configure() *ConfigBuilder {
	return &ConfigBuilder{conf: defaultConfig()}
}

// Addr sets a custom address for the mock HTTP server.
func (cb *ConfigBuilder) Addr(addr string) *ConfigBuilder {
	cb.conf.Addr = addr
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
func (cb *ConfigBuilder) CORS(options ...*CORSConfig) *ConfigBuilder {
	if len(options) > 0 {
		cb.conf.CORS = options[0]
	} else {
		cb.conf.CORS = _defaultCORSConfig
	}

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
func (cb *ConfigBuilder) Parameters(params reply.Params) *ConfigBuilder {
	cb.conf.Parameters = params
	return cb
}

// Files sets a custom Glob patterns to load mock from the file system.
// Defaults to [testdata/*.mock.json, testdata/*.mock.yaml].
func (cb *ConfigBuilder) Files(patterns ...string) *ConfigBuilder {
	cb.conf.Files = patterns
	return cb
}

// Debug allows users to set a function that will be called on unexpected errors.
// This is to help debugging.
func (cb *ConfigBuilder) Debug(debug Debug) *ConfigBuilder {
	cb.conf.Debug = debug
	return cb
}

// Apply builds a new Config with previously configured values.
func (cb *ConfigBuilder) Apply(conf *Config) {
	cb.conf.Apply(conf)
}

// --
// Config Functions
// --

// WithAddr configures the server address.
func WithAddr(addr string) Configurer {
	return configFunc(func(c *Config) { c.Addr = addr })
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
func WithCORS(opts *CORSConfig) Configurer {
	return configFunc(func(c *Config) { c.CORS = opts })
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
func WithParams(params reply.Params) Configurer {
	return configFunc(func(c *Config) { c.Parameters = params })
}

// WithFiles configures directories to search for local mocks.
// Pass a list of glob patterns supported by Go Standard Library.
// This method keeps the default mock filename pattern, [testdata/*mock.json].
// to overwrite the default mock filename pattern, use WithNewFiles.
func WithFiles(patterns ...string) Configurer {
	return configFunc(func(c *Config) { c.Files = append(c.Files, patterns...) })
}

// WithNewFiles configures directories to search for local mocks,
// overwriting the default internal mock filename pattern.
// Pass a list of glob patterns supported by Go Standard Library.
// Use WithFiles to keep the default internal pattern.
func WithNewFiles(patterns ...string) Configurer {
	return configFunc(func(c *Config) { c.Files = patterns })
}

// WithDebug configures a Debug function.
func WithDebug(d Debug) Configurer {
	return configFunc(func(c *Config) { c.Debug = d })
}
