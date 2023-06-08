package httpd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/httpd/cors"
	"github.com/vitorsalgado/mocha/v3/lib"
)

type LogVerbosity int8

const (
	// LogBasic logs only informative messages, without too many details.
	LogBasic LogVerbosity = iota
	LogHeader
	LogBody
)

func (l LogVerbosity) String() string {
	switch l {
	case LogBasic:
		return "basic"
	case LogHeader:
		return "header"
	case LogBody:
		return "body"
	default:
		return "none"
	}
}

type LogLevel = int8

const (
	LogLevelDebug    = LogLevel(zerolog.DebugLevel)
	LogLevelInfo     = LogLevel(zerolog.InfoLevel)
	LogLevelWarn     = LogLevel(zerolog.WarnLevel)
	LogLevelError    = LogLevel(zerolog.ErrorLevel)
	LogLevelNone     = LogLevel(zerolog.NoLevel)
	LogLevelDisabled = LogLevel(zerolog.Disabled)
)

// Defaults
var (
	// ConfigMockFilePattern is the default filename glob pattern to search for local mock files.
	ConfigMockFilePattern = []string{"testdata/_mocks/*mock.yaml", "testdata/_mocks/*mock.yml"}
)

// Configurer lets users configure the Mock API.
type Configurer interface {
	// Apply applies a configuration.
	Apply(conf *Config) error
}

// Config holds HTTPMockApp mock server configurations.
type Config struct {
	// Name sets a name to the mock server.
	// Adds more context for when you have more mocks APIs configured.
	Name string

	// Addr defines a custom server address.
	Addr string

	// UseHTTP2 enables HTTP/2
	UseHTTP2 bool

	// RootDir defines the root directory where the server will start looking for configurations and mocks.
	// Defaults to the current execution path.
	RootDir string

	// TLSConfig defines custom TLS configurations.
	// This is only used when server is started using HTTPMockApp.StartTLS or HTTPMockApp.MustStartTLS.
	// If TLSConfig is set, TLSCertificateFs and TLSKeyFs options will be ignored.
	TLSConfig *tls.Config

	// TLSCertificates is the TLS certificate that will be used by the server
	// All built-in HTTP clients will use this certificate too.
	TLSCertificates []tls.Certificate

	// TLSClientCAs is the client certificate that will be used by the server.
	// All built-in HTTP clients will use it in the RootCAs config.
	TLSClientCAs *x509.CertPool

	// RequestWasNotMatchedStatusCode defines the status code that should be used when
	// an HTTP request doesn't match with any mock.
	// Defaults to 418 (I'm a teapot).
	RequestWasNotMatchedStatusCode int

	// RequestBodyParsers defines request body parsers to be executed before core parsers.
	RequestBodyParsers []RequestBodyParser

	// Middlewares defines a list of custom middlewares that will be
	// set after panic recovery and before mock handler.
	Middlewares []func(http.Handler) http.Handler

	// CORS defines CORS configurations.
	CORS *cors.Config

	// Server defines a custom mock HTTP server.
	Server Server

	// HandlerDecorator provide means to configure a custom HTTP handler
	// while leveraging the default mock handler.
	HandlerDecorator func(handler http.Handler) http.Handler

	// Parameters sets a custom reply parameters store.
	Parameters lib.Params

	// MockFileSearchPatterns configures glob patterns to load mock from the file system.
	MockFileSearchPatterns []string

	// Loaders configures additional loaders.
	Loaders []Loader

	// Proxy configures the mock server as a proxy.
	Proxy *ProxyConfig

	// Record configures Mock Request/Stub recording.
	// Needs to be used with Proxy.
	Record *RecordConfig

	// MockFileHandlers sets custom Mock file handlers for a server instance.
	MockFileHandlers []MockFileHandler

	// TemplateEngine sets a custom template engine.
	TemplateEngine lib.TemplateEngine

	// TemplateFunctions sets custom template functions for the built-in template engine.
	TemplateFunctions template.FuncMap

	// HTTPClientFactory builds an *http.Client that will be used by internal features, like ProxyReply.
	// If none is set, a default one will be used.
	HTTPClientFactory func() *http.Client

	// Logger lets users define a custom logger.
	// If none is provided, a default one will be set.
	Logger *zerolog.Logger

	// LogVerbosity defines the verbosity of the logs.
	LogVerbosity LogVerbosity

	// LogLevel sets the level of the default logger.
	LogLevel LogLevel

	// LogPretty enable/disable pretty logging.
	// This only works with the default zerolog.Logger.
	// If you are setting a custom logger, you need to set this by yourself.
	// Defaults to true.
	LogPretty bool

	// LogBodyMaxSize sets the max size of the response body to be logged.
	// By default, response bodies will be logged entirely.
	LogBodyMaxSize int64

	// UseDescriptiveLogger enable/disable the use of a more descriptive logger for HTTP request matching lifecycle.
	// This is useful, specially for console mode usage, to understand the details of an HTTP request and
	// why a match did not occur.
	// If true, The Logger options will be ignored for the HTTP request matching.
	UseDescriptiveLogger bool

	// Colors enable/disable terminal colors for the descriptive logger.
	// Defaults to true.
	Colors bool

	// PostActions configures PostActions by name.
	// They can later be referenced by name during mocking phase.
	PostActions map[string]PostAction

	// HeaderNamesToRedact configures HTTP header names to be redacted in logs.
	// Useful to avoid logging sensitive information.
	HeaderNamesToRedact map[string]struct{}

	// CLI Only Options

	// UseHTTPS defines that the mock server should use HTTPS.
	// This is only used running the command-line version.
	// To start an HTTPS server from code, call StartTLS() or MustStartTLS() from application instance.
	UseHTTPS bool

	// Forward configures a forward proxy for matched requests.
	Forward *forwardConfig
}

// forwardConfig configures a forward proxy for matched requests.
// Only for CLI.
type forwardConfig struct {
	Target               string
	Headers              http.Header
	ProxyHeaders         http.Header
	ProxyHeadersToRemove []string
	TrimPrefix           string
	TrimSuffix           string
	SSLVerify            bool
}

// Apply copies of the current Config struct values to the given Config parameter.
// It allows the Config struct to be used as a Configurer.
func (c *Config) Apply(conf *Config) error {
	conf.Name = c.Name
	conf.Addr = c.Addr
	conf.RootDir = c.RootDir
	conf.RequestWasNotMatchedStatusCode = c.RequestWasNotMatchedStatusCode
	conf.RequestBodyParsers = c.RequestBodyParsers
	conf.Middlewares = c.Middlewares
	conf.CORS = c.CORS
	conf.Server = c.Server
	conf.HandlerDecorator = c.HandlerDecorator
	conf.Parameters = c.Parameters
	conf.MockFileSearchPatterns = c.MockFileSearchPatterns
	conf.Loaders = c.Loaders
	conf.Proxy = c.Proxy
	conf.Record = c.Record
	conf.Forward = c.Forward
	conf.UseHTTPS = c.UseHTTPS
	conf.MockFileHandlers = c.MockFileHandlers
	conf.TemplateEngine = c.TemplateEngine
	conf.TemplateFunctions = c.TemplateFunctions
	conf.HTTPClientFactory = c.HTTPClientFactory
	conf.UseDescriptiveLogger = c.UseDescriptiveLogger
	conf.Logger = c.Logger
	conf.LogPretty = c.LogPretty
	conf.LogVerbosity = c.LogVerbosity
	conf.LogLevel = c.LogLevel
	conf.LogBodyMaxSize = c.LogBodyMaxSize
	conf.HeaderNamesToRedact = c.HeaderNamesToRedact
	conf.TLSConfig = c.TLSConfig
	conf.TLSCertificates = c.TLSCertificates
	conf.TLSClientCAs = c.TLSClientCAs
	conf.UseHTTP2 = c.UseHTTP2
	conf.PostActions = c.PostActions

	return nil
}

// ConfigBuilder lets users create a Config using a fluent API.
type ConfigBuilder struct {
	conf                   *Config
	recorderConf           []RecordConfigurer
	proxyConf              []ProxyConfigurer
	tlsCertificateFs       string
	tlsKeyFs               string
	tlsClientCertificateFs string
}

func defaultConfig() *Config {
	return &Config{
		RequestWasNotMatchedStatusCode: StatusNoMatch,
		MockFileSearchPatterns:         ConfigMockFilePattern,
		RequestBodyParsers:             make([]RequestBodyParser, 0),
		Middlewares:                    make([]func(http.Handler) http.Handler, 0),
		Loaders:                        make([]Loader, 0),
		MockFileHandlers:               make([]MockFileHandler, 0),
		PostActions:                    make(map[string]PostAction),
		HeaderNamesToRedact:            make(map[string]struct{}),
		UseDescriptiveLogger:           false,
		LogPretty:                      true,
		LogLevel:                       LogLevelInfo,
		LogVerbosity:                   LogHeader,
		Colors:                         true,
		UseHTTP2:                       true,
	}
}

// Setup initialize a new ConfigBuilder.
// Entrypoint to start a new custom configuration for HTTPMockApp mock servers.
func Setup() *ConfigBuilder {
	return &ConfigBuilder{conf: defaultConfig()}
}

// Name sets a name to the mock server.
func (cb *ConfigBuilder) Name(name string) *ConfigBuilder {
	cb.conf.Name = name
	return cb
}

// Addr sets a custom address for the mock HTTP server.
func (cb *ConfigBuilder) Addr(addr string) *ConfigBuilder {
	cb.conf.Addr = addr
	return cb
}

// Port sets a custom port to the mock HTTP server.
// If none is provided, a random port number will be used.
func (cb *ConfigBuilder) Port(port int) *ConfigBuilder {
	cb.conf.Addr = ":" + strconv.FormatInt(int64(port), 10)
	return cb
}

// RootDir defines the root directory where the server will start looking for configurations and mocks.
// Defaults to the current execution path.
func (cb *ConfigBuilder) RootDir(rootDir string) *ConfigBuilder {
	cb.conf.RootDir = rootDir
	return cb
}

// MockNotFoundStatusCode defines the status code to be used when no mock matches with an HTTP request.
func (cb *ConfigBuilder) MockNotFoundStatusCode(code int) *ConfigBuilder {
	cb.conf.RequestWasNotMatchedStatusCode = code
	return cb
}

// RequestBodyParsers adds a custom list of RequestBodyParsers.
func (cb *ConfigBuilder) RequestBodyParsers(bp ...RequestBodyParser) *ConfigBuilder {
	cb.conf.RequestBodyParsers = append(cb.conf.RequestBodyParsers, bp...)
	return cb
}

// Middlewares adds custom middlewares to the mock server.
// Use this to add custom request interceptors.
func (cb *ConfigBuilder) Middlewares(fn ...func(handler http.Handler) http.Handler) *ConfigBuilder {
	cb.conf.Middlewares = append(cb.conf.Middlewares, fn...)
	return cb
}

// CORS configures CORS for the mock server.
func (cb *ConfigBuilder) CORS(options ...cors.Configurer) *ConfigBuilder {
	opts := cors.DefaultConfig
	for _, option := range options {
		option.Apply(&opts)
	}

	cb.conf.CORS = &opts

	return cb
}

// Server configures a custom HTTP mock Server.
func (cb *ConfigBuilder) Server(srv Server) *ConfigBuilder {
	cb.conf.Server = srv
	return cb
}

// HandlerDecorator configures a custom HTTP handler using the default mock handler.
func (cb *ConfigBuilder) HandlerDecorator(fn func(handler http.Handler) http.Handler) *ConfigBuilder {
	cb.conf.HandlerDecorator = fn
	return cb
}

// Logger sets a custom zerolog.Logger.
func (cb *ConfigBuilder) Logger(l *zerolog.Logger) *ConfigBuilder {
	cb.conf.Logger = l
	return cb
}

// LogVerbosity sets the verbosity of informative logs.
// Defaults to LogVerbose.
func (cb *ConfigBuilder) LogVerbosity(l LogVerbosity) *ConfigBuilder {
	cb.conf.LogVerbosity = l
	return cb
}

// LogLevel sets the level of the zerolog.Logger default implementation.
func (cb *ConfigBuilder) LogLevel(l LogLevel) *ConfigBuilder {
	cb.conf.LogLevel = l
	return cb
}

// LogPretty enable/disable pretty logging.
// This only works with the default zerolog.Logger.
// If you are setting a custom logger, you need to set this by yourself.
// Defaults to true.
func (cb *ConfigBuilder) LogPretty(v bool) *ConfigBuilder {
	cb.conf.LogPretty = v
	return cb
}

// LogBodyMaxSize sets the max size of the response body to be logged.
// By default, response bodies will be logged entirely.
func (cb *ConfigBuilder) LogBodyMaxSize(max int64) *ConfigBuilder {
	cb.conf.LogBodyMaxSize = max
	return cb
}

// UseDescriptiveLogger enable/disable the use of a more descriptive logger for HTTP request matching lifecycle.
// This is useful, specially for console mode usage, to understand the details of an HTTP request and
// why a match did not occur.
// If true, The Logger options will be ignored for the HTTP request matching.
func (cb *ConfigBuilder) UseDescriptiveLogger() *ConfigBuilder {
	cb.conf.UseDescriptiveLogger = true
	return cb
}

// RedactHeader redacts the given HTTP header names in logs.
func (cb *ConfigBuilder) RedactHeader(names ...string) *ConfigBuilder {
	for _, name := range names {
		cb.conf.HeaderNamesToRedact[strings.ToLower(name)] = struct{}{}
	}

	return cb
}

// Parameters sets a custom reply parameters store.
func (cb *ConfigBuilder) Parameters(params lib.Params) *ConfigBuilder {
	cb.conf.Parameters = params
	return cb
}

// MockFilePatterns sets custom Glob patterns to load mock from the file system.
// Defaults to [testdata/*.mock.json, testdata/*.mock.yaml].
func (cb *ConfigBuilder) MockFilePatterns(patterns ...string) *ConfigBuilder {
	cb.conf.MockFileSearchPatterns = patterns

	return cb
}

// Loader configures an additional Loader.
func (cb *ConfigBuilder) Loader(loader Loader) *ConfigBuilder {
	cb.conf.Loaders = append(cb.conf.Loaders, loader)
	return cb
}

// Proxy configures the mock server as a proxy server.
// Non-Matched requests will be routed based on the proxy configuration.
func (cb *ConfigBuilder) Proxy(options ...ProxyConfigurer) *ConfigBuilder {
	if len(options) == 0 {
		c := _defaultProxyConfig
		cb.proxyConf = []ProxyConfigurer{&c}

		return cb
	}

	cb.proxyConf = options

	return cb
}

// Record configures recording.
func (cb *ConfigBuilder) Record(options ...RecordConfigurer) *ConfigBuilder {
	if len(options) == 0 {
		cb.recorderConf = []RecordConfigurer{defaultRecordConfig()}
		return cb
	}

	cb.recorderConf = options

	return cb
}

// MockFileHandlers sets MockFileHandler implementations.
func (cb *ConfigBuilder) MockFileHandlers(handlers ...MockFileHandler) *ConfigBuilder {
	cb.conf.MockFileHandlers = append(cb.conf.MockFileHandlers, handlers...)
	return cb
}

// TemplateEngine sets the TemplateEngine to be used by all components.
// Defaults to a built-in implementation based on Go Templates.
func (cb *ConfigBuilder) TemplateEngine(te lib.TemplateEngine) *ConfigBuilder {
	cb.conf.TemplateEngine = te
	return cb
}

// TemplateEngineFunctions sets custom functions to be used in templates.
// This only works with the built-in TemplateEngine implementation.
// Custom template engine implementations must provide their own mean to set custom functions.
func (cb *ConfigBuilder) TemplateEngineFunctions(fm template.FuncMap) *ConfigBuilder {
	cb.conf.TemplateFunctions = fm
	return cb
}

// HTTPClient sets a custom http.Client factory.
// Internal components that require an HTTP client will use this factory,
// instead of using the default implementation.
func (cb *ConfigBuilder) HTTPClient(f func() *http.Client) *ConfigBuilder {
	cb.conf.HTTPClientFactory = f
	return cb
}

// TLSConfig defines custom TLS configurations.
// This is only used when server is started using HTTPMockApp.StartTLS or Mocha.MustStartTLS.
// If TLSConfig is set, TLSCertificateFs and TLSKeyFs options will be ignored.
func (cb *ConfigBuilder) TLSConfig(c *tls.Config) *ConfigBuilder {
	cb.conf.TLSConfig = c
	return cb
}

// TLSCertKeyPair sets a custom public/private key pair.
// If none is provided, default values will be used when starting the server with HTTPMockApp.StartTLS or HTTPMockApp.MustStartTLS.
// If TLSConfig is set, this option will be ignored.
func (cb *ConfigBuilder) TLSCertKeyPair(certFile, keyFile string) *ConfigBuilder {
	cb.tlsCertificateFs = certFile
	cb.tlsKeyFs = keyFile
	return cb
}

// TLSMutual defines the set of root certificate authorities that the server will use.
// If TLSConfig is set, this option will be ignored.
func (cb *ConfigBuilder) TLSMutual(certFile, keyFile, clientCert string) *ConfigBuilder {
	cb.tlsClientCertificateFs = clientCert
	cb.tlsCertificateFs = certFile
	cb.tlsKeyFs = keyFile
	return cb
}

func (cb *ConfigBuilder) PostAction(name string, action PostAction) *ConfigBuilder {
	cb.conf.PostActions[name] = action
	return cb
}

// Apply builds a new Config with previously configured values.
func (cb *ConfigBuilder) Apply(conf *Config) error {
	if len(cb.recorderConf) > 0 {
		recordConfig := defaultRecordConfig()
		for _, option := range cb.recorderConf {
			err := option.Apply(recordConfig)
			if err != nil {
				return err
			}
		}

		cb.conf.Record = recordConfig
	}

	if len(cb.proxyConf) > 0 {
		proxyConfig := _defaultProxyConfig
		for _, option := range cb.proxyConf {
			err := option.Apply(&proxyConfig)
			if err != nil {
				return err
			}
		}

		cb.conf.Proxy = &proxyConfig
	}

	if len(cb.tlsCertificateFs) == 0 && len(cb.tlsKeyFs) == 0 {
		return cb.conf.Apply(conf)
	}

	err := applyTLS(cb.conf, cb.tlsCertificateFs, cb.tlsKeyFs, cb.tlsClientCertificateFs)
	if err != nil {
		return err
	}

	return cb.conf.Apply(conf)
}

func applyTLS(c *Config, certFile, keyFile, clientCertFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("config: failed to load certificate key pair from file=%s key=%s. %w", certFile, keyFile, err)
	}

	if len(clientCertFile) > 0 {
		cc, err := os.ReadFile(filepath.Clean(clientCertFile))
		if err != nil {
			return fmt.Errorf("config: failed to load client certificate from %s. %w", clientCertFile, err)
		}

		cp := x509.NewCertPool()
		cp.AppendCertsFromPEM(cc)

		c.TLSClientCAs = cp
	}

	c.TLSCertificates = []tls.Certificate{cert}

	return nil
}
