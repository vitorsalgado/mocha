package mocha

import (
	"context"
	"net/http"

	"github.com/vitorsalgado/mocha/feat/cors"
)

type LogVerbosity int

// Log verbosity enum
const (
	LogSilently LogVerbosity = iota
	LogVerbose
)

type (
	// Config holds Mocha mock server configurations.
	Config struct {
		// Context to be used internally by Mocha.
		Context context.Context

		// Addr defines a custom server address.
		Addr string

		// BodyParsers defines request body parsers to be executed before core parsers.
		BodyParsers []RequestBodyParser

		// Middlewares defines a list of custom middlewares that will be
		// set after panic recover and before mock handler.
		Middlewares []func(http.Handler) http.Handler

		// CORS defines CORS configurations.
		CORS cors.Config

		// Server defines a custom mock HTTP server.
		Server Server

		// LogVerbosity defines the level of logs
		LogVerbosity LogVerbosity

		corsEnabled bool
	}

	// Configurer is Config builder,
	// Use this to build Mocha options, instead of creating a new Config struct manually.
	Configurer struct {
		conf Config
	}
)

var configDefault = Configure().LogVerbosity(LogVerbose).Build()

// Configure inits a new Configurer.
// Entrypoint to start a new custom configuration for Mocha mock servers.
func Configure() *Configurer {
	return &Configurer{conf: Config{
		LogVerbosity: LogVerbose,
		BodyParsers:  make([]RequestBodyParser, 0),
		Middlewares:  make([]func(http.Handler) http.Handler, 0)}}
}

// Context sets a custom context.
func (cb *Configurer) Context(context context.Context) *Configurer {
	cb.conf.Context = context
	return cb
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

// CORS configures CORS for the mock server.
func (cb *Configurer) CORS(options ...cors.Config) *Configurer {
	if len(options) > 0 {
		cb.conf.CORS = options[0]
	} else {
		cb.conf.CORS = cors.ConfigDefault
	}

	cb.conf.corsEnabled = true

	return cb
}

// Server configures a custom HTTP mock Server.
func (cb *Configurer) Server(srv Server) *Configurer {
	cb.conf.Server = srv
	return cb
}

// LogVerbosity configure the verbosity of informative logs.
// Defaults to LogVerbose.
func (cb *Configurer) LogVerbosity(l LogVerbosity) *Configurer {
	cb.conf.LogVerbosity = l
	return cb
}

// Build builds a new Config with previously configured values.
func (cb *Configurer) Build() Config {
	return cb.conf
}
