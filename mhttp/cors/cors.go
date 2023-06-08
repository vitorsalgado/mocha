package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vitorsalgado/mocha/v3/mhttpv"
)

// Config represents the possible options to configure CORS.
type Config struct {
	AllowedOrigin     string
	AllowCredentials  bool
	AllowedMethods    string
	AllowedHeaders    string
	ExposeHeaders     string
	MaxAge            int
	SuccessStatusCode int
}

// Configurer lets users configure CORS.
type Configurer interface {
	Apply(opts *Config)
}

// Apply allows CORSConfig to be used as a CORSConfigurer
func (c *Config) Apply(opts *Config) {
	opts.AllowedOrigin = c.AllowedOrigin
	opts.AllowCredentials = c.AllowCredentials
	opts.AllowedMethods = c.AllowedMethods
	opts.AllowedHeaders = c.AllowedHeaders
	opts.ExposeHeaders = c.ExposeHeaders
	opts.MaxAge = c.MaxAge
	opts.SuccessStatusCode = c.SuccessStatusCode
}

var DefaultConfig = Config{
	AllowedOrigin: "*",
	AllowedMethods: strings.Join([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodHead,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}, ","),
	AllowedHeaders:    "",
	AllowCredentials:  false,
	ExposeHeaders:     "",
	MaxAge:            0,
	SuccessStatusCode: http.StatusNoContent,
}

// ConfigBuilder facilitates building CORS options.
type ConfigBuilder struct {
	options *Config
	origins []string
}

// CORS initializes a CORSConfig builder for a fluent configuration.
func CORS() *ConfigBuilder {
	return &ConfigBuilder{
		origins: make([]string, 0),
		options: &Config{SuccessStatusCode: http.StatusNoContent}}
}

// SuccessStatusCode sets a custom status code returned on CORS Options request.
// If none is specified, the default status code is http.StatusNoContent.
func (b *ConfigBuilder) SuccessStatusCode(code int) *ConfigBuilder {
	b.options.SuccessStatusCode = code
	return b
}

// MaxAge sets CORS max age.
func (b *ConfigBuilder) MaxAge(maxAge int) *ConfigBuilder {
	b.options.MaxAge = maxAge
	return b
}

// AllowOrigin sets the allowed origins.
func (b *ConfigBuilder) AllowOrigin(origin ...string) *ConfigBuilder {
	b.origins = append(b.origins, origin...)
	return b
}

// AllowCredentials sets "Access-Control-Allow-Credentials" misc.Header
func (b *ConfigBuilder) AllowCredentials(allow bool) *ConfigBuilder {
	b.options.AllowCredentials = allow
	return b
}

// ExposeHeaders sets "Access-Control-Expose-Header" misc.Header
func (b *ConfigBuilder) ExposeHeaders(headers ...string) *ConfigBuilder {
	b.options.ExposeHeaders = strings.Join(headers, ",")
	return b
}

// AllowedHeaders sets the allowed headers.
// It will set the header "Access-Control-Allow-Header".
func (b *ConfigBuilder) AllowedHeaders(headers ...string) *ConfigBuilder {
	b.options.AllowedHeaders = strings.Join(headers, ",")
	return b
}

// AllowMethods sets the allowed HTTP methods.
// The header "Access-Control-Allow-Methods" will be used.
func (b *ConfigBuilder) AllowMethods(methods ...string) *ConfigBuilder {
	b.options.AllowedMethods = strings.Join(methods, ",")
	return b
}

// Apply builds CORS configurations based on previous settings set via the builder.
func (b *ConfigBuilder) Apply(opts *Config) {
	b.build().Apply(opts)
}

func (b *ConfigBuilder) build() *Config {
	if len(b.origins) > 0 {
		if len(b.origins) == 1 {
			b.options.AllowedOrigin = b.origins[0]
		} else {
			b.options.AllowedOrigin = strings.Join(b.origins, ",")
		}
	}

	return b.options
}

func New(options *Config) func(http.Handler) http.Handler {
	if options == nil {
		options = &DefaultConfig
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// preflight request
			if r.Method == http.MethodOptions {
				configureOrigin(options, r, w)
				configureCredentials(options, w)
				configureExposedHeaders(options, w)
				configureMethods(options, w)
				configureMaxAge(options, w)
				configureHeaders(options, w, r)

				w.Header().Add(mhttpv.HeaderVary, mhttpv.HeaderAccessControlRequestHeaders)
				w.Header().Add(mhttpv.HeaderContentLength, "0")

				w.WriteHeader(options.SuccessStatusCode)
			} else {
				configureOrigin(options, r, w)
				configureCredentials(options, w)
				configureExposedHeaders(options, w)

				next.ServeHTTP(w, r)
				return
			}
		})
	}
}

func configureHeaders(options *Config, w http.ResponseWriter, r *http.Request) {
	// when allowed headers aren't specified, use values from header access-control-request-headers
	if options.AllowedHeaders != "" {
		w.Header().Add(mhttpv.HeaderAccessControlAllowHeaders, options.AllowedHeaders)
	} else {
		hs := r.Header.Get(mhttpv.HeaderAccessControlRequestHeaders)
		if strings.TrimSpace(hs) != "" {
			w.Header().Add(mhttpv.HeaderAccessControlAllowHeaders, hs)
		}
	}
}

func configureMaxAge(options *Config, w http.ResponseWriter) {
	if options.MaxAge > -1 {
		w.Header().Add(mhttpv.HeaderAccessControlMaxAge, strconv.Itoa(options.MaxAge))
	}
}

func configureMethods(options *Config, w http.ResponseWriter) {
	if len(options.AllowedMethods) > 0 {
		w.Header().Add(mhttpv.HeaderAccessControlAllowMethods, options.AllowedMethods)
	}
}

func configureExposedHeaders(options *Config, w http.ResponseWriter) {
	if options.ExposeHeaders != "" {
		w.Header().Add(mhttpv.HeaderAccessControlExposeHeaders, options.ExposeHeaders)
	}
}

func configureCredentials(options *Config, w http.ResponseWriter) {
	if options.AllowCredentials {
		w.Header().Add(mhttpv.HeaderAccessControlAllowCredentials, "true")
	}
}

func configureOrigin(options *Config, r *http.Request, w http.ResponseWriter) {
	if options.AllowedOrigin == "" {
		return
	}

	origins := strings.Split(options.AllowedOrigin, ",")
	size := len(origins)

	if size == 1 {
		w.Header().Add(mhttpv.HeaderAccessControlAllowOrigin, options.AllowedOrigin)
		w.Header().Add(mhttpv.HeaderVary, mhttpv.HeaderOrigin)
		return
	}

	// received a list of origins
	// will check if the request origin is within the provided array and use it as the allowed origin
	origin := r.Header.Get("origin")
	allowed := false

	for _, o := range origins {
		if origin == o {
			allowed = true
			break
		}
	}

	if allowed {
		w.Header().Add(mhttpv.HeaderAccessControlAllowOrigin, origin)
		w.Header().Add(mhttpv.HeaderVary, mhttpv.HeaderOrigin)
	}
}
