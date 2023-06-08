package mhttp

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/vitorsalgado/mocha/v3/httpd/cors"
)

const (
	// DefaultConfigFileName is the default configuration filename.
	DefaultConfigFileName = ".moairc"
	// EnvPrefix is the standard prefix that should be present in all environment variables used by this application.
	EnvPrefix = "MOAI"
)

var (
	// DefaultConfigDirectories is the default configuration directories used to lookup for a configuration file.
	DefaultConfigDirectories = []string{".", "testdata"}
)

const (
	_kName              = "name"
	_kAddr              = "addr"
	_kRootDir           = "root_dir"
	_kMockFiles         = "mock_files"
	_kNoMatchStatusCode = "no_match_status_code"
	_kConfig            = "config"
	_kUseHTTPS          = "https"
	_kColors            = "colors"

	_kTLS        = "tls"
	_kTLSCert    = "tls.cert"
	_kTLSKey     = "tls.key"
	_kTLSRootCAs = "tls.ca"

	_kLog               = "logger"
	_kLogLevel          = "logger.level"
	_kLogVerbosity      = "logger.verbosity"
	_kLogMaxBodySize    = "logger.max_body_size"
	_kLogPretty         = "logger.pretty"
	_kLogUseDescriptive = "logger.descriptive"

	_kCORS                  = "cors"
	_kCORSAllowedOrigin     = "cors.allowed_origin"
	_kCORSAllowCredentials  = "cors.allow_credentials"
	_kCORSAllowedMethods    = "cors.allowed_methods"
	_kCORSAllowedHeaders    = "cors.allowed_headers"
	_kCORSExposeHeaders     = "cors.expose_headers"
	_kCORSMaxAge            = "cors.max_age"
	_kCORSSuccessStatusCode = "cors.success_status_code"

	_kProxy          = "proxy"
	_kProxyVia       = "proxy.via"
	_kProxyTimeout   = "proxy.timeout"
	_kProxySSLVerify = "proxy.ssl_verify"

	_kRecord                = "record"
	_kRecordRequestHeaders  = "record.request_headers"
	_kRecordResponseHeaders = "record.response_headers"
	_kRecordSaveDir         = "record.save_dir"
	_kRecordSaveBodyToFile  = "record.save_body_to_file"
	_kRecordFileType        = "record.file_type"

	_kForward                   = "forward"
	_kForwardTarget             = "forward.target"
	_kForwardHeaders            = "forward.headers"
	_kForwardProxyHeaders       = "forward.proxy_headers"
	_kForwardRemoveProxyHeaders = "forward.remove_proxy_headers"
	_kForwardTrimPrefix         = "forward.trim_prefix"
	_kForwardTrimSuffix         = "forward.trim_suffix"
	_kForwardSSLVerify          = "forward.ssl_verify"
)

var _ Configurer = (*localConfigurer)(nil)

var _once sync.Once

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

// UseConfig enables lookup for local configuration files using standard naming conventions.
// Supported extensions (json|toml|yaml|yml|properties|props|prop|hcl|tfvars|dotenv|env|ini)".
// If only the filename is provided, it must contain the full path and extension to the configuration.
func UseConfig(filename string, paths ...string) Configurer {
	return &localConfigurer{filename: filename, paths: paths}
}

// Apply applies configurations using viper.Viper.
func (c *localConfigurer) Apply(conf *Config) (err error) {
	v := viper.New()

	_once.Do(func() {
		err = bindFlags(v)
	})

	if err != nil {
		return fmt.Errorf("config: error binding command line flags.\n%w", err)
	}

	err = bindEnv(v)
	if err != nil {
		return fmt.Errorf("config: error binding environment variables.\n%w", err)
	}

	// use config filename from flags, if present,
	// before binding local file configurations.
	filename := v.GetString(_kConfig)
	if filename != "" {
		c.filename = filename
	}

	err = bindFile(v, c.filename, c.paths)
	if err != nil {
		return fmt.Errorf(
			"config: error binding configuration from file=%s, lookup_paths=%v.\n%w",
			c.filename,
			c.paths,
			err,
		)
	}

	err = c.apply(conf, v)
	if err != nil {
		return fmt.Errorf("config: error applying configuration from %s.\n%w", c.filename, err)
	}

	return nil
}

func (c *localConfigurer) apply(conf *Config, vi *viper.Viper) (err error) {
	vi.SetDefault(_kName, conf.Name)
	vi.SetDefault(_kAddr, conf.Addr)
	vi.SetDefault(_kLogLevel, conf.LogLevel)
	vi.SetDefault(_kConfig, DefaultConfigFileName)
	vi.SetDefault(_kLogVerbosity, int8(conf.LogVerbosity))
	vi.SetDefault(_kLogPretty, conf.LogPretty)
	vi.SetDefault(_kLogUseDescriptive, conf.UseDescriptiveLogger)
	vi.SetDefault(_kLogMaxBodySize, conf.LogBodyMaxSize)
	vi.SetDefault(_kColors, conf.Colors)
	vi.SetDefault(_kMockFiles, conf.MockFileSearchPatterns)
	vi.SetDefault(_kRootDir, conf.RootDir)
	vi.SetDefault(_kUseHTTPS, conf.UseHTTPS)
	vi.SetDefault(_kNoMatchStatusCode, conf.RequestWasNotMatchedStatusCode)

	conf.Name = vi.GetString(_kName)
	conf.Addr = vi.GetString(_kAddr)
	conf.RootDir = vi.GetString(_kRootDir)
	conf.MockFileSearchPatterns = vi.GetStringSlice(_kMockFiles)
	conf.Colors = vi.GetBool(_kColors)
	conf.RequestWasNotMatchedStatusCode = vi.GetInt(_kNoMatchStatusCode)

	if vi.IsSet(_kTLS) {
		certFile := vi.GetString(_kTLSCert)
		keyFile := vi.GetString(_kTLSKey)
		clientCertFile := vi.GetString(_kTLSRootCAs)

		if len(certFile) == 0 || len(keyFile) == 0 {
			return fmt.Errorf("both cert and key filename are required. got cert=%s key=%s", certFile, keyFile)
		}

		err = applyTLS(conf, certFile, keyFile, clientCertFile)
		if err != nil {
			return err
		}
	}

	if vi.IsSet(_kLog) {
		conf.LogLevel = LogLevel(vi.GetInt(_kLogLevel))
		conf.LogVerbosity = LogVerbosity(vi.GetInt(_kLogVerbosity))
		conf.LogBodyMaxSize = vi.GetInt64(_kLogMaxBodySize)
		conf.LogPretty = vi.GetBool(_kLogPretty)
		conf.UseDescriptiveLogger = vi.GetBool(_kLogUseDescriptive)
	}

	if vi.IsSet(_kCORS) {
		def := cors.DefaultConfig

		vi.SetDefault(_kCORSAllowedOrigin, def.AllowedOrigin)
		vi.SetDefault(_kCORSAllowCredentials, def.AllowCredentials)
		vi.SetDefault(_kCORSAllowedMethods, def.AllowedMethods)
		vi.SetDefault(_kCORSAllowedHeaders, def.AllowedHeaders)
		vi.SetDefault(_kCORSExposeHeaders, def.ExposeHeaders)
		vi.SetDefault(_kCORSMaxAge, def.MaxAge)
		vi.SetDefault(_kCORSSuccessStatusCode, def.SuccessStatusCode)

		conf.CORS = &cors.Config{
			AllowedOrigin:     vi.GetString(_kCORSAllowedOrigin),
			AllowCredentials:  vi.GetBool(_kCORSAllowCredentials),
			AllowedMethods:    vi.GetString(_kCORSAllowedMethods),
			AllowedHeaders:    vi.GetString(_kCORSAllowedHeaders),
			ExposeHeaders:     vi.GetString(_kCORSExposeHeaders),
			MaxAge:            vi.GetInt(_kCORSMaxAge),
			SuccessStatusCode: vi.GetInt(_kCORSSuccessStatusCode),
		}
	}

	if vi.IsSet(_kProxy) {
		vi.SetDefault(_kProxyTimeout, _defaultProxyConfig.Timeout.Milliseconds())
		vi.SetDefault(_kProxySSLVerify, true)

		vv := vi.Get(_kProxy)
		switch vv.(type) {
		case bool:
			conf.Proxy = &ProxyConfig{Timeout: _defaultProxyConfig.Timeout}
		case map[string]any:
			conf.Proxy = &ProxyConfig{
				Via:       vi.GetString(_kProxyVia),
				Timeout:   time.Duration(vi.GetInt64(_kProxyTimeout)),
				SSLVerify: vi.GetBool(_kProxySSLVerify),
			}
		default:
			return errors.New(`field "proxy" has an unsupported type. supported types are: object, bool`)
		}
	}

	if vi.IsSet(_kRecord) {
		defRec := defaultRecordConfig()

		vi.SetDefault(_kRecordRequestHeaders, defRec.RequestHeaders)
		vi.SetDefault(_kRecordResponseHeaders, defRec.ResponseHeaders)
		vi.SetDefault(_kRecordSaveDir, defRec.SaveDir)
		vi.SetDefault(_kRecordSaveBodyToFile, defRec.SaveResponseBodyToFile)
		vi.SetDefault(_kRecordFileType, defRec.SaveFileType)

		conf.Record = &RecordConfig{
			RequestHeaders:         vi.GetStringSlice(_kRecordRequestHeaders),
			ResponseHeaders:        vi.GetStringSlice(_kRecordResponseHeaders),
			SaveDir:                vi.GetString(_kRecordSaveDir),
			SaveResponseBodyToFile: vi.GetBool(_kRecordSaveBodyToFile),
			SaveFileType:           vi.GetString(_kRecordFileType),
		}
	}

	if vi.IsSet(_kForward) || vi.IsSet(_kForwardTarget) {
		target := vi.GetString(_kForwardTarget)
		if target == "" {
			return errors.New(`when specifying a "forward" configuration, the field "forward.target" is required`)
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			return fmt.Errorf(`field "forward.target" must contain a valid URL. error=%w`, err)
		}

		h := make(http.Header, len(vi.GetStringMapString(_kForwardHeaders)))
		for k, v := range vi.GetStringMapString(_kForwardHeaders) {
			h.Add(k, v)
		}

		hp := make(http.Header, len(vi.GetStringMapString(_kForwardProxyHeaders)))
		for k, v := range vi.GetStringMapString(_kForwardProxyHeaders) {
			hp.Add(k, v)
		}

		conf.Forward = &forwardConfig{
			Target:               targetURL.String(),
			Headers:              h,
			ProxyHeaders:         hp,
			ProxyHeadersToRemove: vi.GetStringSlice(_kForwardRemoveProxyHeaders),
			TrimPrefix:           vi.GetString(_kForwardTrimPrefix),
			TrimSuffix:           vi.GetString(_kForwardTrimSuffix),
			SSLVerify:            vi.GetBool(_kForwardSSLVerify),
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
	cl := pflag.CommandLine
	if cl.Parsed() {
		return v.BindPFlags(cl)
	}

	normalizer := cl.GetNormalizeFunc()
	cl.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		name = strings.ToLower(string(normalizer(f, name)))
		name = strings.ReplaceAll(name, ".", "-")
		name = strings.ReplaceAll(name, "_", "-")

		return pflag.NormalizedName(name)
	})

	cl.StringP(_kConfig, "c", "", "Use a custom configuration file. Commandline arguments take precedence over file values.")
	cl.StringP(_kName, "n", "", "Mock server name.")
	cl.String(_kRootDir, "", "Root directory to start looking for configurations and mocks.")
	cl.StringP(_kAddr, "a", "",
		"Server address. Usage: <host>:<port>, :<port>, <port>. If no value is set, it will use a localhost with a random port.")
	cl.Bool(_kColors, true, "Enabled/disable colors for descriptive logger only.")
	cl.StringSlice(_kMockFiles, nil, "Mock files search glob patterns. E.g. testdata/*mock.yaml.")
	cl.Bool(_kUseHTTPS, false, "Enabled HTTPS.")

	cl.String(_kTLSCert, "", "TLS certificate file.")
	cl.String(_kTLSKey, "", "TLS private key file.")
	cl.String(_kTLSRootCAs, "", "TLS trusted CA certificates filenames.")

	cl.Int8(_kLogLevel, LogLevelInfo, "Log level.")
	cl.Int8(_kLogVerbosity, int8(LogHeader), "Log verbosity.")
	cl.Bool(_kLogPretty, true, "Pretty logging.")
	cl.Bool(_kLogUseDescriptive, false, "Use descriptive logger.")
	cl.Int64(_kLogMaxBodySize, 0, "Max body size to be logged.")

	cl.Bool(_kRecord, false, "Enabled mock recording.")
	cl.String(_kRecordSaveDir, "", "Recorded mocks directory.")
	cl.Bool(_kRecordSaveBodyToFile, false, "Save recorded body to separate file, using \"body_file\" field.")
	cl.StringSlice(_kRecordRequestHeaders, nil, "Request headers to be recorded.")
	cl.StringSlice(_kRecordResponseHeaders, nil, "Response headers to be recorded.")
	cl.String(_kRecordFileType, "", "Recorded mock file format. E.g. json,yaml")

	cl.Bool(_kProxy, false, "Enabled reverse proxy.")
	cl.String(_kProxyVia, "", "Route proxy requests to another proxy.")
	cl.Int64(_kProxyTimeout, _defaultProxyConfig.Timeout.Milliseconds(), "Proxy timeout in milliseconds.")
	cl.Bool(_kProxySSLVerify, false, "Verify SSL certificates.")

	cl.Bool(_kCORS, false, "Enabled CORS.")

	cl.Bool(_kForward, false, "")
	cl.String(_kForwardTarget, "", "Forward requests to the given URL.")
	cl.Bool(_kForwardSSLVerify, false, "SSL verify")

	err := cl.Parse(os.Args)
	if err != nil {
		return err
	}

	return v.BindPFlags(cl)
}

func bindEnv(v *viper.Viper) error {
	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return nil
}
