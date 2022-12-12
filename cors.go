package mocha

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vitorsalgado/mocha/v3/internal/header"
)

// CORSConfig represents the possible options to configure corsMid for the mock server.
type CORSConfig struct {
	AllowedOrigin     string
	AllowCredentials  bool
	AllowedMethods    string
	AllowedHeaders    string
	ExposeHeaders     string
	MaxAge            int
	SuccessStatusCode int
}

// CORSConfigurer lets users configure CORS.
type CORSConfigurer interface {
	Apply(opts *CORSConfig)
}

// Apply allows CORSConfig to be used as a CORSConfigurer
func (c *CORSConfig) Apply(opts *CORSConfig) {
	opts.AllowedOrigin = c.AllowedOrigin
	opts.AllowCredentials = c.AllowCredentials
	opts.AllowedMethods = c.AllowedMethods
	opts.AllowedHeaders = c.AllowedHeaders
	opts.ExposeHeaders = c.ExposeHeaders
	opts.MaxAge = c.MaxAge
	opts.SuccessStatusCode = c.SuccessStatusCode
}

var _defaultCORSConfig = CORSConfig{
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

// CORSConfigBuilder facilitates building corsMid options.
type CORSConfigBuilder struct {
	options *CORSConfig
	origins []string
}

// CORS inits a CORSConfig builder for a fluent configuration.
func CORS() *CORSConfigBuilder {
	return &CORSConfigBuilder{
		origins: make([]string, 0),
		options: &CORSConfig{SuccessStatusCode: http.StatusNoContent}}
}

// SuccessStatusCode sets a custom status code returned on corsMid Options request.
// If none is specified, the default status code is http.StatusNoContent.
func (b *CORSConfigBuilder) SuccessStatusCode(code int) *CORSConfigBuilder {
	b.options.SuccessStatusCode = code
	return b
}

// MaxAge sets corsMid max age.
func (b *CORSConfigBuilder) MaxAge(maxAge int) *CORSConfigBuilder {
	b.options.MaxAge = maxAge
	return b
}

// AllowOrigin sets allowed origins.
func (b *CORSConfigBuilder) AllowOrigin(origin ...string) *CORSConfigBuilder {
	b.origins = append(b.origins, origin...)
	return b
}

// AllowCredentials sets "Access-Control-Allow-Credentials" header.
func (b *CORSConfigBuilder) AllowCredentials(allow bool) *CORSConfigBuilder {
	b.options.AllowCredentials = allow
	return b
}

// ExposeHeaders sets "Access-Control-Expose-Header" header.
func (b *CORSConfigBuilder) ExposeHeaders(headers ...string) *CORSConfigBuilder {
	b.options.ExposeHeaders = strings.Join(headers, ",")
	return b
}

// AllowedHeaders sets allowed headers.
// It will set the header "Access-Control-Allow-Header".
func (b *CORSConfigBuilder) AllowedHeaders(headers ...string) *CORSConfigBuilder {
	b.options.AllowedHeaders = strings.Join(headers, ",")
	return b
}

// AllowMethods sets the allowed HTTP methods.
// The header "Access-Control-Allow-Methods" will be used.
func (b *CORSConfigBuilder) AllowMethods(methods ...string) *CORSConfigBuilder {
	b.options.AllowedMethods = strings.Join(methods, ",")
	return b
}

// Apply builds CORS configurations based on previous settings set via the builder.
func (b *CORSConfigBuilder) Apply(opts *CORSConfig) {
	b.build().Apply(opts)
}

func (b *CORSConfigBuilder) build() *CORSConfig {
	if len(b.origins) > 0 {
		if len(b.origins) == 1 {
			b.options.AllowedOrigin = b.origins[0]
		} else {
			b.options.AllowedOrigin = strings.Join(b.origins, ",")
		}
	}

	return b.options
}

func corsMid(options *CORSConfig) func(http.Handler) http.Handler {
	if options == nil {
		options = &_defaultCORSConfig
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

				w.Header().Add(header.Vary, header.AccessControlRequestHeaders)
				w.Header().Add(header.ContentLength, "0")

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

func configureHeaders(options *CORSConfig, w http.ResponseWriter, r *http.Request) {
	// when allowed headers aren't specified, use values from header access-control-request-headers
	if options.AllowedHeaders != "" {
		w.Header().Add(header.AccessControlAllowHeaders, options.AllowedHeaders)
	} else {
		hs := r.Header.Get(header.AccessControlRequestHeaders)
		if strings.TrimSpace(hs) != "" {
			w.Header().Add(header.AccessControlAllowHeaders, hs)
		}
	}
}

func configureMaxAge(options *CORSConfig, w http.ResponseWriter) {
	if options.MaxAge > -1 {
		w.Header().Add(header.AccessControlMaxAge, strconv.Itoa(options.MaxAge))
	}
}

func configureMethods(options *CORSConfig, w http.ResponseWriter) {
	if len(options.AllowedMethods) > 0 {
		w.Header().Add(header.AccessControlAllowMethods, options.AllowedMethods)
	}
}

func configureExposedHeaders(options *CORSConfig, w http.ResponseWriter) {
	if options.ExposeHeaders != "" {
		w.Header().Add(header.AccessControlExposeHeaders, options.ExposeHeaders)
	}
}

func configureCredentials(options *CORSConfig, w http.ResponseWriter) {
	if options.AllowCredentials {
		w.Header().Add(header.AccessControlAllowCredentials, "true")
	}
}

func configureOrigin(options *CORSConfig, r *http.Request, w http.ResponseWriter) {
	if options.AllowedOrigin == "" {
		return
	}

	origins := strings.Split(options.AllowedOrigin, ",")
	size := len(origins)

	if size == 1 {
		w.Header().Add(header.AccessControlAllowOrigin, options.AllowedOrigin)
		w.Header().Add(header.Vary, header.Origin)
		return
	}

	// received a list of origins
	// will check if request origin is within the provided array and use it as the allowed origin
	origin := r.Header.Get("origin")
	allowed := false

	for _, o := range origins {
		if origin == o {
			allowed = true
			break
		}
	}

	if allowed {
		w.Header().Add(header.AccessControlAllowOrigin, origin)
		w.Header().Add(header.Vary, header.Origin)
	}
}
