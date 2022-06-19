package mocha

type (
	Config struct {
		BodyParsers []BodyParser
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

func (cb *ConfigBuilder) Build() *Config {
	return cb.conf
}
