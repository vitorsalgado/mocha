package mocha

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/reply"
)

type LogLevel int

const (
	LogSilently LogLevel = iota
	LogInfo
	LogVerbose
)

// Configurer lets users configure the Mock API.
type Configurer interface {
	// Apply applies a configuration.
	Apply(conf *Config)
}

// Config holds Mocha mock server configurations.
type Config struct {
	// Addr defines a custom server address.
	Addr string

	// BodyParsers defines request body parsers to be executed before core parsers.
	BodyParsers []RequestBodyParser

	// Middlewares defines a list of custom middlewares that will be
	// set after panic recover and before mock handler.
	Middlewares []func(http.Handler) http.Handler

	// CORS defines CORS configurations.
	CORS *CORSConfig

	// Server defines a custom mock HTTP server.
	Server Server

	// Handler provides a mean to configure a custom HTTP handler
	// while leveraging the default mock handler.
	Handler func(handler http.Handler) http.Handler

	// LogLevel defines the level of logs
	LogLevel LogLevel

	// Parameters sets a custom reply parameters store.
	Parameters reply.Params

	// FileMockPatterns configures glob patterns to load mock from the file system.
	FileMockPatterns []string

	Debug Debug
}

// Apply copies the current Config struct values to the given Config parameter.
// It allows the Config struct to be used as a Configurer.
func (c *Config) Apply(conf *Config) {
	conf.Addr = c.Addr
	conf.BodyParsers = c.BodyParsers
	conf.Middlewares = c.Middlewares
	conf.CORS = c.CORS
	conf.Server = c.Server
	conf.Handler = c.Handler
	conf.LogLevel = c.LogLevel
	conf.Parameters = c.Parameters
	conf.FileMockPatterns = c.FileMockPatterns
	conf.Debug = c.Debug
}

// ConfigBuilder lets users create a Config with a fluent API.
type ConfigBuilder struct {
	conf *Config
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
	cb.conf.BodyParsers = append(cb.conf.BodyParsers, bp...)
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
	cb.conf.Handler = fn
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

// MockFilePatterns sets a custom Glob patterns to load mock from the file system.
// Defaults to [testdata/*.mock.json, testdata/*.mock.yaml].
func (cb *ConfigBuilder) MockFilePatterns(patterns ...string) *ConfigBuilder {
	cb.conf.FileMockPatterns = patterns
	return cb
}

func (cb *ConfigBuilder) Debug(debug Debug) *ConfigBuilder {
	cb.conf.Debug = debug
	return cb
}

// Apply builds a new Config with previously configured values.
func (cb *ConfigBuilder) Apply(conf *Config) {
	cb.conf.Apply(conf)
}

func defaultConfig() *Config {
	return &Config{
		LogLevel:    LogVerbose,
		BodyParsers: make([]RequestBodyParser, 0),
		Middlewares: make([]func(http.Handler) http.Handler, 0)}
}
