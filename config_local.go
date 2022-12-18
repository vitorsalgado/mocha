package mocha

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// DefaultConfigFileName is the default configuration filename.
	DefaultConfigFileName = ".moairc"
)

var (
	// DefaultConfigDirectories is the default configuration directories used lookup for a configuration file.
	DefaultConfigDirectories = []string{".", "testdata"}
)

const (
	_fieldName        = "name"
	_fieldAddr        = "addr"
	_fieldLogLevel    = "log_level"
	_fieldDirectories = "directories"

	_fieldCORS                  = "cors"
	_fieldCORSAllowedOrigin     = "cors.allowed_origin"
	_fieldCORSAllowCredentials  = "cors.allow_credentials"
	_fieldCORSAllowedMethods    = "cors.allowed_methods"
	_fieldCORSAllowedHeaders    = "cors.allowed_headers"
	_fieldCORSExposeHeaders     = "cors.expose_headers"
	_fieldCORSMaxAge            = "cors.max_age"
	_fieldCORSSuccessStatusCode = "cors.success_status_code"

	_fieldProxy        = "proxy"
	_fieldProxyTarget  = "proxy.target"
	_fieldProxyTimeout = "proxy.timeout"

	_fieldRecord                = "record"
	_fieldRecordRequestHeaders  = "record.request_headers"
	_fieldRecordResponseHeaders = "record.response_headers"
	_fieldRecordSave            = "record.save"
	_fieldRecordSaveDir         = "record.save_dir"
	_fieldRecordSaveBodyToFile  = "record.save_body_file"
)

var _ Configurer = (*localConfigurer)(nil)

type localConfigurer struct {
	filename string
	paths    []string
}

// WithLocal enables lookup for local configuration files using standard naming conventions.
// It will look up for a file named ".moairc.(json|toml|yaml|yml|properties|props|prop|hcl|tfvars|dotenv|env|ini)"
// in the directories "." and "testdata".
func WithLocal() Configurer {
	return &localConfigurer{filename: DefaultConfigFileName, paths: DefaultConfigDirectories}
}

// WithLocalCustomized enables lookup for local configuration file using standard naming conventions.
// Supported extensions (json|toml|yaml|yml|properties|props|prop|hcl|tfvars|dotenv|env|ini)".
// If only the filename is provided, it must contain the full path and extension to the configuration.
func WithLocalCustomized(filename string, paths []string) Configurer {
	return &localConfigurer{filename: filename, paths: paths}
}

// Apply applies configurations using viper.Viper.
func (c *localConfigurer) Apply(conf *Config) error {
	v := viper.New()

	// when no paths are set, use the configuration filename.
	if len(c.paths) == 0 {
		v.SetConfigFile(c.filename)
	} else {
		filename := c.filename

		for _, fileType := range viper.SupportedExts {
			if strings.HasSuffix(filename, fileType) {
				filename = strings.TrimSuffix(filename, "."+fileType)
				break
			}
		}

		v.SetConfigName(filename)

		for _, p := range c.paths {
			v.AddConfigPath(p)
		}
	}

	v.SetDefault(_fieldName, conf.Name)
	v.SetDefault(_fieldAddr, conf.Addr)
	v.SetDefault(_fieldLogLevel, int(conf.LogLevel))
	v.SetDefault(_fieldDirectories, conf.Directories)

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return err
	}

	conf.Name = v.GetString(_fieldName)
	conf.Addr = v.GetString(_fieldAddr)
	conf.LogLevel = LogLevel(v.GetInt(_fieldLogLevel))
	conf.Directories = v.GetStringSlice(_fieldDirectories)

	if v.IsSet(_fieldCORS) {
		v.SetDefault(_fieldCORSAllowedOrigin, _defaultCORSConfig.AllowedOrigin)
		v.SetDefault(_fieldCORSAllowCredentials, _defaultCORSConfig.AllowCredentials)
		v.SetDefault(_fieldCORSAllowedMethods, _defaultCORSConfig.AllowedMethods)
		v.SetDefault(_fieldCORSAllowedHeaders, _defaultCORSConfig.AllowedHeaders)
		v.SetDefault(_fieldCORSExposeHeaders, _defaultCORSConfig.ExposeHeaders)
		v.SetDefault(_fieldCORSMaxAge, _defaultCORSConfig.MaxAge)
		v.SetDefault(_fieldCORSSuccessStatusCode, _defaultCORSConfig.SuccessStatusCode)

		conf.CORS = &CORSConfig{
			AllowedOrigin:     v.GetString(_fieldCORSAllowedOrigin),
			AllowCredentials:  v.GetBool(_fieldCORSAllowCredentials),
			AllowedMethods:    v.GetString(_fieldCORSAllowedMethods),
			AllowedHeaders:    v.GetString(_fieldCORSAllowedHeaders),
			ExposeHeaders:     v.GetString(_fieldCORSExposeHeaders),
			MaxAge:            v.GetInt(_fieldCORSMaxAge),
			SuccessStatusCode: v.GetInt(_fieldCORSSuccessStatusCode),
		}
	}

	if v.IsSet(_fieldProxy) {
		v.SetDefault(_fieldProxyTimeout, _defaultProxyConfig.Timeout.Milliseconds())

		conf.Proxy = &ProxyConfig{
			Target:  v.GetString(_fieldProxyTarget),
			Timeout: time.Duration(v.GetInt64(_fieldProxyTimeout)),
		}
	}

	if v.IsSet(_fieldRecord) {
		v.SetDefault(_fieldRecordRequestHeaders, _defaultRecordConfig.RequestHeaders)
		v.SetDefault(_fieldRecordResponseHeaders, _defaultRecordConfig.ResponseHeaders)
		v.SetDefault(_fieldRecordSave, _defaultRecordConfig.Save)
		v.SetDefault(_fieldRecordSaveDir, _defaultRecordConfig.SaveDir)
		v.SetDefault(_fieldRecordSaveBodyToFile, _defaultRecordConfig.SaveBodyToFile)

		conf.Record = &RecordConfig{
			RequestHeaders:  v.GetStringSlice(_fieldRecordRequestHeaders),
			ResponseHeaders: v.GetStringSlice(_fieldRecordResponseHeaders),
			Save:            v.GetBool(_fieldRecordSave),
			SaveDir:         v.GetString(_fieldRecordSaveDir),
			SaveBodyToFile:  v.GetBool(_fieldRecordSaveBodyToFile),
		}
	}

	return nil
}
