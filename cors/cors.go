// Package cors implements functions to enable CORS support for the mock server.
package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vitorsalgado/mocha/util/headers"
)

// Options represents the possible options to configure CORS for the mock server.
type Options struct {
	AllowedOrigin     string
	AllowCredentials  string
	AllowedMethods    string
	AllowedHeaders    string
	ExposeHeaders     string
	MaxAge            int
	SuccessStatusCode int
}

// CORS returns a http.Handler that will be used to handle CORS requests.
// To build options more easily, use the options' builder cors.Configure().
// Example:
//	cors.CORS(cors.Configure().AllowOrigin(""))
func CORS(options Options) func(http.Handler) http.Handler {
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

				w.Header().Add(headers.Vary, headers.AccessControlRequestHeaders)
				w.Header().Add(headers.ContentLength, "0")

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

func configureHeaders(options Options, w http.ResponseWriter, r *http.Request) {
	// when allowed headers aren't specified, use values from header access-control-request-headers
	if options.AllowedHeaders != "" {
		w.Header().Add(headers.AccessControlAllowHeaders, options.AllowedHeaders)
	} else {

		hs := r.Header.Get(headers.AccessControlRequestHeaders)
		if strings.TrimSpace(hs) != "" {
			w.Header().Add(headers.AccessControlAllowHeaders, hs)
		}
	}
}

func configureMaxAge(options Options, w http.ResponseWriter) {
	if options.MaxAge > -1 {
		w.Header().Add(headers.AccessControlMaxAge, strconv.Itoa(options.MaxAge))
	}
}

func configureMethods(options Options, w http.ResponseWriter) {
	if len(options.AllowedMethods) > 0 {
		w.Header().Add(headers.AccessControlAllowMethods, options.AllowedMethods)
	}
}

func configureExposedHeaders(options Options, w http.ResponseWriter) {
	if options.ExposeHeaders != "" {
		w.Header().Add(headers.AccessControlExposeHeaders, options.ExposeHeaders)
	}
}

func configureCredentials(options Options, w http.ResponseWriter) {
	if options.AllowCredentials != "" {
		w.Header().Add(headers.AccessControlAllowCredentials, options.AllowCredentials)
	}
}

func configureOrigin(options Options, r *http.Request, w http.ResponseWriter) {
	if options.AllowedOrigin == "" {
		return
	}

	origins := strings.Split(options.AllowedOrigin, ",")
	size := len(origins)

	if size == 1 {
		w.Header().Add(headers.AccessControlAllowOrigin, options.AllowedOrigin)
		w.Header().Add(headers.Vary, headers.Origin)
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
		w.Header().Add(headers.AccessControlAllowOrigin, origin)
		w.Header().Add(headers.Vary, headers.Origin)
	}
}
