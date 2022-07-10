package mocha

import (
	"context"
	"net/http"

	"github.com/vitorsalgado/mocha/cors"
)

type (
	// Config holds Mocha mock server configurations.
	Config struct {
		Context     context.Context
		Addr        string
		BodyParsers []RequestBodyParser
		Middlewares []func(http.Handler) http.Handler
		CORS        cors.Options
	}

	// Configurer is Config builder,
	// Use this to build Mocha options, instead of creating a new Config struct manually.
	Configurer struct {
		conf Config
	}
)

var configDefault = Configure().Build()

// Configure inits a new Configurer.
// Entrypoint to start a new custom configuration for Mocha mock servers.
func Configure() *Configurer {
	return &Configurer{conf: Config{BodyParsers: make([]RequestBodyParser, 0)}}
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

// RequestBodyParser adds a custom list of RequestBodyParser.
func (cb *Configurer) RequestBodyParser(bp ...RequestBodyParser) *Configurer {
	cb.conf.BodyParsers = append(cb.conf.BodyParsers, bp...)
	return cb
}

// Middlewares adds custom middlewares to the mock server.
// Use this to add custom request interceptors.
func (cb *Configurer) Middlewares(fn ...func(next http.Handler) http.Handler) *Configurer {
	cb.conf.Middlewares = append(cb.conf.Middlewares, fn...)
	return cb
}

// CORS configures CORS for the mock server.
func (cb *Configurer) CORS(options *cors.OptionsBuilder) *Configurer {
	cb.conf.CORS = options.Build()
	return cb
}

// Build builds a new Config with previously configured values.
func (cb *Configurer) Build() Config {
	return cb.conf
}
