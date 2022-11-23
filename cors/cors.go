// Package cors implements functions to enable cross-origin resource sharing support for the mock server.
package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vitorsalgado/mocha/v3/internal/headerx"
)

// Config represents the possible options to configure New for the mock server.
type Config struct {
	AllowedOrigin     string
	AllowCredentials  bool
	AllowedMethods    string
	AllowedHeaders    string
	ExposeHeaders     string
	MaxAge            int
	SuccessStatusCode int
}

// ConfigDefault is the default config.
var ConfigDefault = Config{
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

// New returns a http.Handler that will be used to handle New requests.
// To build options more easily, use the options' builder cors.Configure().
// Example:
//
//	cors.New(cors.Configure().AllowOrigin(""))
func New(options Config) func(http.Handler) http.Handler {
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

				w.Header().Add(headerx.Vary, headerx.AccessControlRequestHeaders)
				w.Header().Add(headerx.ContentLength, "0")

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

func configureHeaders(options Config, w http.ResponseWriter, r *http.Request) {
	// when allowed headers aren't specified, use values from header access-control-request-headers
	if options.AllowedHeaders != "" {
		w.Header().Add(headerx.AccessControlAllowHeaders, options.AllowedHeaders)
	} else {
		hs := r.Header.Get(headerx.AccessControlRequestHeaders)
		if strings.TrimSpace(hs) != "" {
			w.Header().Add(headerx.AccessControlAllowHeaders, hs)
		}
	}
}

func configureMaxAge(options Config, w http.ResponseWriter) {
	if options.MaxAge > -1 {
		w.Header().Add(headerx.AccessControlMaxAge, strconv.Itoa(options.MaxAge))
	}
}

func configureMethods(options Config, w http.ResponseWriter) {
	if len(options.AllowedMethods) > 0 {
		w.Header().Add(headerx.AccessControlAllowMethods, options.AllowedMethods)
	}
}

func configureExposedHeaders(options Config, w http.ResponseWriter) {
	if options.ExposeHeaders != "" {
		w.Header().Add(headerx.AccessControlExposeHeaders, options.ExposeHeaders)
	}
}

func configureCredentials(options Config, w http.ResponseWriter) {
	if options.AllowCredentials {
		w.Header().Add(headerx.AccessControlAllowCredentials, "true")
	}
}

func configureOrigin(options Config, r *http.Request, w http.ResponseWriter) {
	if options.AllowedOrigin == "" {
		return
	}

	origins := strings.Split(options.AllowedOrigin, ",")
	size := len(origins)

	if size == 1 {
		w.Header().Add(headerx.AccessControlAllowOrigin, options.AllowedOrigin)
		w.Header().Add(headerx.Vary, headerx.Origin)
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
		w.Header().Add(headerx.AccessControlAllowOrigin, origin)
		w.Header().Add(headerx.Vary, headerx.Origin)
	}
}
