package dzgrpc

import (
	"google.golang.org/grpc"
)

type Config struct {
	Service       any
	ServiceDesc   *grpc.ServiceDesc
	ServerOptions []grpc.ServerOption
	Addr          string
}

func (c *Config) Apply(conf *Config) error {
	conf.Service = c.Service
	conf.ServiceDesc = c.ServiceDesc
	conf.ServerOptions = c.ServerOptions
	conf.Addr = c.Addr

	return nil
}

func defaultConfig() *Config {
	return &Config{}
}

type ConfigBuilder struct {
	conf *Config
}

func Setup() *ConfigBuilder {
	return &ConfigBuilder{conf: defaultConfig()}
}

func (cb *ConfigBuilder) Service(sd *grpc.ServiceDesc, service any) *ConfigBuilder {
	cb.conf.ServiceDesc = sd
	cb.conf.Service = service

	return cb
}

func (cb *ConfigBuilder) ServerOptions(options ...grpc.ServerOption) *ConfigBuilder {
	cb.conf.ServerOptions = options
	return cb
}

// Addr sets a custom address for the mock GRPC server.
func (cb *ConfigBuilder) Addr(addr string) *ConfigBuilder {
	cb.conf.Addr = addr
	return cb
}

func (cb *ConfigBuilder) Apply(conf *Config) error {
	return cb.conf.Apply(conf)
}
