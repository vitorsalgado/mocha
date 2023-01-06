package mocha

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// DefaultConfigFileName is the default configuration filename.
	DefaultConfigFileName = ".moairc"
	EnvPrefix             = "MOAI"
)

var (
	// DefaultConfigDirectories is the default configuration directories used lookup for a configuration file.
	DefaultConfigDirectories = []string{".", "testdata"}
)

const (
	_kName        = "name"
	_kNameP       = "n"
	_kAddr        = "addr"
	_kAddrP       = "a"
	_kLogLevel    = "log_level"
	_kDirectories = "directories"
	_kGlob        = "glob"
	_kGlobP       = "g"

	_kConfig  = "config"
	_kConfigP = "c"

	_kUseHTTPS  = "https"
	_kUseHTTPsP = "s"

	_kCORS                  = "cors"
	_kCORSAllowedOrigin     = "cors.allowed_origin"
	_kCORSAllowCredentials  = "cors.allow_credentials"
	_kCORSAllowedMethods    = "cors.allowed_methods"
	_kCORSAllowedHeaders    = "cors.allowed_headers"
	_kCORSExposeHeaders     = "cors.expose_headers"
	_kCORSMaxAge            = "cors.max_age"
	_kCORSSuccessStatusCode = "cors.success_status_code"

	_kProxy        = "proxy"
	_kProxyP       = "p"
	_kProxyVia     = "proxy.proxy_via"
	_kProxyTimeout = "proxy.timeout"

	_kRecord                = "record"
	_kRecordP               = "r"
	_kRecordRequestHeaders  = "record.request_headers"
	_kRecordResponseHeaders = "record.response_headers"
	_kRecordSave            = "record.save"
	_kRecordSaveDir         = "record.save_dir"
	_kRecordSaveBodyToFile  = "record.save_body_file"

	_kForward                   = "forward"
	_kForwardTarget             = "forward.target"
	_kForwardHeaders            = "forward.headers"
	_kForwardProxyHeaders       = "forward.proxy_headers"
	_kForwardRemoveProxyHeaders = "forward.remove_proxy_headers"
	_kForwardTrimPrefix         = "forward.trim_prefix"
	_kForwardTrimSuffix         = "forward.trim_suffix"
)

var _ Configurer = (*localConfigurer)(nil)

type localConfigurer struct {
	filename string
	paths    []string
}

// UseLocals enables lookup for local configuration files using standard naming conventions.
// It will look up for a file named ".moairc.(json|toml|yaml|yml|properties|props|prop|hcl|tfvars|dotenv|env|ini)"
// in the directories "." and "testdata".
func UseLocals() Configurer {
	return &localConfigurer{filename: DefaultConfigFileName, paths: DefaultConfigDirectories}
}

// UseLocalsWith enables lookup for local configuration file using standard naming conventions.
// Supported extensions (json|toml|yaml|yml|properties|props|prop|hcl|tfvars|dotenv|env|ini)".
// If only the filename is provided, it must contain the full path and extension to the configuration.
func UseLocalsWith(filename string, paths []string) Configurer {
	return &localConfigurer{filename: filename, paths: paths}
}

// Apply applies configurations using viper.Viper.
func (c *localConfigurer) Apply(conf *Config) (err error) {
	v := viper.New()

	err = bindFlags(v)
	if err != nil {
		return err
	}

	err = bindEnv(v)
	if err != nil {
		return err
	}

	// use config filename from flags, if present,
	// before binding local file configurations.
	filename := v.GetString(_kConfig)
	if filename != "" {
		c.filename = filename
	}

	err = bindFile(v, c.filename, c.paths)
	if err != nil {
		return err
	}

	v.SetDefault(_kName, conf.Name)
	v.SetDefault(_kAddr, conf.Addr)
	v.SetDefault(_kLogLevel, int(conf.LogLevel))
	v.SetDefault(_kDirectories, conf.Directories)

	conf.Name = v.GetString(_kName)
	conf.Addr = v.GetString(_kAddr)
	conf.LogLevel = LogLevel(v.GetInt(_kLogLevel))
	conf.Directories = v.GetStringSlice(_kDirectories)

	if v.IsSet(_kCORS) {
		v.SetDefault(_kCORSAllowedOrigin, _defaultCORSConfig.AllowedOrigin)
		v.SetDefault(_kCORSAllowCredentials, _defaultCORSConfig.AllowCredentials)
		v.SetDefault(_kCORSAllowedMethods, _defaultCORSConfig.AllowedMethods)
		v.SetDefault(_kCORSAllowedHeaders, _defaultCORSConfig.AllowedHeaders)
		v.SetDefault(_kCORSExposeHeaders, _defaultCORSConfig.ExposeHeaders)
		v.SetDefault(_kCORSMaxAge, _defaultCORSConfig.MaxAge)
		v.SetDefault(_kCORSSuccessStatusCode, _defaultCORSConfig.SuccessStatusCode)

		conf.CORS = &CORSConfig{
			AllowedOrigin:     v.GetString(_kCORSAllowedOrigin),
			AllowCredentials:  v.GetBool(_kCORSAllowCredentials),
			AllowedMethods:    v.GetString(_kCORSAllowedMethods),
			AllowedHeaders:    v.GetString(_kCORSAllowedHeaders),
			ExposeHeaders:     v.GetString(_kCORSExposeHeaders),
			MaxAge:            v.GetInt(_kCORSMaxAge),
			SuccessStatusCode: v.GetInt(_kCORSSuccessStatusCode),
		}
	}

	if v.IsSet(_kProxy) {
		v.SetDefault(_kProxyTimeout, _defaultProxyConfig.Timeout.Milliseconds())

		vv := v.Get(_kProxy)
		switch vv.(type) {
		case bool:
			conf.Proxy = &ProxyConfig{Timeout: _defaultProxyConfig.Timeout}
		case map[string]any:
			conf.Proxy = &ProxyConfig{
				ProxyVia: v.GetString(_kProxyVia),
				Timeout:  time.Duration(v.GetInt64(_kProxyTimeout)),
			}
		default:
			return errors.New(`field "proxy" has an unknown type. supported type are: object, bool`)
		}
	}

	if v.IsSet(_kRecord) {
		defRec := defaultRecordConfig()

		v.SetDefault(_kRecordRequestHeaders, defRec.RequestHeaders)
		v.SetDefault(_kRecordResponseHeaders, defRec.ResponseHeaders)
		v.SetDefault(_kRecordSave, defRec.Save)
		v.SetDefault(_kRecordSaveDir, defRec.SaveDir)
		v.SetDefault(_kRecordSaveBodyToFile, defRec.SaveBodyToFile)

		conf.Record = &RecordConfig{
			RequestHeaders:  v.GetStringSlice(_kRecordRequestHeaders),
			ResponseHeaders: v.GetStringSlice(_kRecordResponseHeaders),
			Save:            v.GetBool(_kRecordSave),
			SaveDir:         v.GetString(_kRecordSaveDir),
			SaveBodyToFile:  v.GetBool(_kRecordSaveBodyToFile),
		}
	}

	if v.IsSet(_kForward) || v.IsSet(_kForwardTarget) {
		target := v.GetString(_kForwardTarget)
		if target == "" {
			return errors.New(`when specifying a "forward" configuration, the field "forward.target" is required`)
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			return fmt.Errorf(`field "forward.target" must contain a valid URL. %w`, err)
		}

		h := make(http.Header, len(v.GetStringMapString(_kForwardHeaders)))
		for k, v := range v.GetStringMapString(_kForwardHeaders) {
			h.Add(k, v)
		}

		hp := make(http.Header, len(v.GetStringMapString(_kForwardProxyHeaders)))
		for k, v := range v.GetStringMapString(_kForwardProxyHeaders) {
			hp.Add(k, v)
		}

		conf.Forward = &ForwardConfig{
			Target:               targetURL.String(),
			Headers:              h,
			ProxyHeaders:         hp,
			ProxyHeadersToRemove: v.GetStringSlice(_kForwardRemoveProxyHeaders),
			TrimPrefix:           v.GetString(_kForwardTrimPrefix),
			TrimSuffix:           v.GetString(_kForwardTrimSuffix),
		}
	}

	return nil
}

func bindFile(v *viper.Viper, filename string, paths []string) error {
	// when no paths are set, use the configuration filename.
	if len(paths) == 0 {
		v.SetConfigFile(filename)
	} else {
		for _, fileType := range viper.SupportedExts {
			if strings.HasSuffix(filename, fileType) {
				filename = strings.TrimSuffix(filename, "."+fileType)
				break
			}
		}

		v.SetConfigName(filename)

		for _, p := range paths {
			v.AddConfigPath(p)
		}
	}

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func bindFlags(v *viper.Viper) error {
	pp := pflag.CommandLine

	normalizer := pp.GetNormalizeFunc()

	pp.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		result := normalizer(f, name)
		name = strings.ReplaceAll(strings.ToLower(string(result)), ".", "-")

		return pflag.NormalizedName(name)
	})

	pp.StringP(_kConfig, _kConfigP, "", "Use a custom configuration file. Commandline arguments take precedence over file values.")
	pp.StringP(_kName, _kNameP, "", "Mock server name.")
	pp.BoolP(_kUseHTTPS, _kUseHTTPsP, false, "Enable HTTPS.")
	pp.StringP(_kAddr, _kAddrP, "", "Server address. Usage: <host>:<port>, :<port>, <port>. If no value is set, it will an auto generated one.")
	pp.Int(_kLogLevel, int(LogVerbose), "Verbose logs")
	pp.StringSliceP(_kGlob, _kGlobP, nil, "Mock search glob patterns. Example: testdata/*mock.json,testdata/*mock.yaml")

	pp.BoolP(_kRecord, _kRecordP, false, "Enable mock recording.")
	pp.String(_kRecordSaveDir, "testdata/", "Recorded mocks directory.")
	pp.Bool(_kRecordSaveBodyToFile, false, "Save recorded body to separate file, using \"body_file\" field.")
	pp.StringSlice(_kRecordRequestHeaders, nil, "Request headers to be recorded.")
	pp.StringSlice(_kRecordResponseHeaders, nil, "Response headers to be recorded.")

	pp.BoolP(_kProxy, _kProxyP, false, "Enable reverse proxy.")
	pp.String(_kProxyVia, "", "Route proxy requests to another proxy.")
	pp.Int64(_kProxyTimeout, _defaultProxyConfig.Timeout.Milliseconds(), "Proxy timeout in milliseconds.")

	pp.Bool(_kCORS, false, "Enable CORS.")

	pp.Bool(_kForward, false, "")
	pp.String(_kForwardTarget, "", "Forward requests to the given URL.")

	v.RegisterAlias(_kDirectories, _kGlob)

	err := pp.Parse(os.Args)
	if err != nil {
		return err
	}

	err = v.BindPFlags(pp)
	if err != nil {
		return err
	}

	return nil
}

func bindEnv(v *viper.Viper) error {
	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}
