package mocha

import (
	"time"

	"github.com/spf13/viper"
)

const (
	ConfigFileName = ".moai"
)

var ConfigPaths = []string{".", "testdata"}

const (
	_ckCORSAllowedOrigin     = "cors.allowed_origin"
	_ckCORSAllowCredentials  = "cors.allow_credentials"
	_ckCORSAllowedMethods    = "cors.allowed_methods"
	_ckCORSAllowedHeaders    = "cors.allowed_headers"
	_ckCORSExposeHeaders     = "cors.expose_headers"
	_ckCORSMaxAge            = "cors.max_age"
	_ckCORSSuccessStatusCode = "cors.success_status_code"

	_ckProxyTarget  = "proxy.target"
	_ckProxyTimeout = "proxy.timeout"

	_ckRecordRequestHeaders  = "record.request_headers"
	_ckRecordResponseHeaders = "record.response_headers"
	_ckRecordSave            = "record.save"
	_ckRecordSaveDir         = "record.save_dir"
	_ckRecordSaveBodyToFile  = "record.save_body_file"
)

type builtInConfigurer struct {
	filename string
	paths    []string
}

func BuiltInConfigurer() Configurer {
	return &builtInConfigurer{filename: ConfigFileName, paths: ConfigPaths}
}

func BuiltInConfigurerWith(filename string, paths []string) Configurer {
	return &builtInConfigurer{filename: filename, paths: paths}
}

func (c *builtInConfigurer) Apply(conf *Config) error {
	v := viper.New()
	v.SetConfigName(c.filename)

	for _, p := range c.paths {
		v.AddConfigPath(p)
	}

	v.SetDefault("name", conf.Name)
	v.SetDefault("addr", conf.Addr)
	v.SetDefault("log_level", int(conf.LogLevel))
	v.SetDefault("files", conf.Directories)

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return err
	}

	conf.Name = v.GetString("name")
	conf.Addr = v.GetString("addr")
	conf.LogLevel = LogLevel(v.GetInt("log_level"))
	conf.Directories = v.GetStringSlice("files")

	if v.IsSet("cors") {
		v.SetDefault(_ckCORSAllowedOrigin, _defaultCORSConfig.AllowedOrigin)
		v.SetDefault(_ckCORSAllowCredentials, _defaultCORSConfig.AllowCredentials)
		v.SetDefault(_ckCORSAllowedMethods, _defaultCORSConfig.AllowedMethods)
		v.SetDefault(_ckCORSAllowedHeaders, _defaultCORSConfig.AllowedHeaders)
		v.SetDefault(_ckCORSExposeHeaders, _defaultCORSConfig.ExposeHeaders)
		v.SetDefault(_ckCORSMaxAge, _defaultCORSConfig.MaxAge)
		v.SetDefault(_ckCORSSuccessStatusCode, _defaultCORSConfig.SuccessStatusCode)

		conf.CORS = &CORSConfig{
			AllowedOrigin:     v.GetString(_ckCORSAllowedOrigin),
			AllowCredentials:  v.GetBool(_ckCORSAllowCredentials),
			AllowedMethods:    v.GetString(_ckCORSAllowedMethods),
			AllowedHeaders:    v.GetString(_ckCORSAllowedHeaders),
			ExposeHeaders:     v.GetString(_ckCORSExposeHeaders),
			MaxAge:            v.GetInt(_ckCORSMaxAge),
			SuccessStatusCode: v.GetInt(_ckCORSSuccessStatusCode),
		}
	}

	if v.IsSet("proxy") {
		v.SetDefault(_ckProxyTimeout, _defaultProxyConfig.Timeout.Milliseconds())

		conf.Proxy = &ProxyConfig{
			Target:  v.GetString(_ckProxyTarget),
			Timeout: time.Duration(v.GetInt64(_ckProxyTimeout)),
		}
	}

	if v.IsSet("record") {
		v.SetDefault(_ckRecordRequestHeaders, _defaultRecordConfig.RequestHeaders)
		v.SetDefault(_ckRecordResponseHeaders, _defaultRecordConfig.ResponseHeaders)
		v.SetDefault(_ckRecordSave, _defaultRecordConfig.Save)
		v.SetDefault(_ckRecordSaveDir, _defaultRecordConfig.SaveDir)
		v.SetDefault(_ckRecordSaveBodyToFile, _defaultRecordConfig.SaveBodyToFile)

		conf.Record = &RecordConfig{
			RequestHeaders:  v.GetStringSlice(_ckRecordRequestHeaders),
			ResponseHeaders: v.GetStringSlice(_ckRecordResponseHeaders),
			Save:            v.GetBool(_ckRecordSave),
			SaveDir:         v.GetString(_ckRecordSaveDir),
			SaveBodyToFile:  v.GetBool(_ckRecordSaveBodyToFile),
		}
	}

	return nil
}
