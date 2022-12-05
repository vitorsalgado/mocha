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

	Pattern string

	corsEnabled bool
}

// Configurer is Config builder,
// Use this to build Mocha options, instead of creating a new Config struct manually.
type Configurer struct {
	conf *Config
}

var _configDefault = Configure().
	LogLevel(LogVerbose).
	Build()

// Configure inits a new Configurer.
// Entrypoint to start a new custom configuration for Mocha mock servers.
func Configure() *Configurer {
	return &Configurer{conf: &Config{
		LogLevel:    LogVerbose,
		BodyParsers: make([]RequestBodyParser, 0),
		Middlewares: make([]func(http.Handler) http.Handler, 0)}}
}

// Addr sets a custom address for the mock HTTP server.
func (cb *Configurer) Addr(addr string) *Configurer {
	cb.conf.Addr = addr
	return cb
}

// RequestBodyParsers adds a custom list of RequestBodyParsers.
func (cb *Configurer) RequestBodyParsers(bp ...RequestBodyParser) *Configurer {
	cb.conf.BodyParsers = append(cb.conf.BodyParsers, bp...)
	return cb
}

// Middlewares adds custom middlewares to the mock server.
// Use this to add custom request interceptors.
func (cb *Configurer) Middlewares(fn ...func(handler http.Handler) http.Handler) *Configurer {
	cb.conf.Middlewares = append(cb.conf.Middlewares, fn...)
	return cb
}

// CORS configures Cross Origin Resource Sharing for the mock server.
func (cb *Configurer) CORS(options ...*CORSConfig) *Configurer {
	if len(options) > 0 {
		cb.conf.CORS = options[0]
	} else {
		cb.conf.CORS = _defaultCORSConfig
	}

	cb.conf.corsEnabled = true

	return cb
}

// Server configures a custom HTTP mock Server.
func (cb *Configurer) Server(srv Server) *Configurer {
	cb.conf.Server = srv
	return cb
}

// HandlerDecorator configures a custom HTTP handler using the default mock handler.
func (cb *Configurer) HandlerDecorator(fn func(handler http.Handler) http.Handler) *Configurer {
	cb.conf.Handler = fn
	return cb
}

// LogLevel configure the verbosity of informative logs.
// Defaults to LogVerbose.
func (cb *Configurer) LogLevel(l LogLevel) *Configurer {
	cb.conf.LogLevel = l
	return cb
}

// Parameters sets a custom reply parameters store.
func (cb *Configurer) Parameters(params reply.Params) *Configurer {
	cb.conf.Parameters = params
	return cb
}

func (cb *Configurer) Pattern(pattern string) *Configurer {
	cb.conf.Pattern = pattern
	return cb
}

// Build builds a new Config with previously configured values.
func (cb *Configurer) Build() *Config {
	return cb.conf
}
