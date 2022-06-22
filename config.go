package mocha

import (
	"context"
	"net/http"
)

type (
	Config struct {
		Context     context.Context
		Addr        string
		BodyParsers []BodyParser
		Middlewares []func(http.Handler) http.Handler
		CORS        *CORSOptions
	}

	Conf interface {
		Build() *Config
	}

	ConfigBuilder struct {
		conf *Config
	}
)

func Options() *ConfigBuilder {
	return &ConfigBuilder{conf: &Config{BodyParsers: make([]BodyParser, 0)}}
}

func (cb *ConfigBuilder) Context(context context.Context) *ConfigBuilder {
	cb.conf.Context = context
	return cb
}

func (cb *ConfigBuilder) Addr(addr string) *ConfigBuilder {
	cb.conf.Addr = addr
	return cb
}

func (cb *ConfigBuilder) BodyParsers(bp ...BodyParser) *ConfigBuilder {
	cb.conf.BodyParsers = append(cb.conf.BodyParsers, bp...)
	return cb
}

func (cb *ConfigBuilder) Middlewares(fn ...func(next http.Handler) http.Handler) *ConfigBuilder {
	cb.conf.Middlewares = append(cb.conf.Middlewares, fn...)
	return cb
}

func (cb *ConfigBuilder) CORS(options *CORSOptionsBuilder) *ConfigBuilder {
	cb.conf.CORS = options.Build()
	return cb
}

func (cb *ConfigBuilder) Build() *Config {
	return cb.conf
}
