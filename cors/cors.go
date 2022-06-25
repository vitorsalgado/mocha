package cors

import (
	"github.com/vitorsalgado/mocha/internal/header"
	"net/http"
	"strconv"
	"strings"
)

type Options struct {
	AllowedOrigin     string
	AllowCredentials  string
	AllowedMethods    string
	AllowedHeaders    string
	ExposeHeaders     string
	MaxAge            int
	SuccessStatusCode int
}

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

func configureHeaders(options Options, w http.ResponseWriter, r *http.Request) {
	// when allowed headers aren't specified, use values from header access-control-request-headers
	if options.AllowedHeaders != "" {
		w.Header().Add(header.AccessControlAllowHeaders, options.AllowedHeaders)
	} else {

		headers := r.Header.Get(header.AccessControlRequestHeaders)
		if strings.TrimSpace(headers) != "" {
			w.Header().Add(header.AccessControlAllowHeaders, headers)
		}
	}
}

func configureMaxAge(options Options, w http.ResponseWriter) {
	if options.MaxAge > -1 {
		w.Header().Add(header.AccessControlMaxAge, strconv.Itoa(options.MaxAge))
	}
}

func configureMethods(options Options, w http.ResponseWriter) {
	if len(options.AllowedMethods) > 0 {
		w.Header().Add(header.AccessControlAllowMethods, options.AllowedMethods)
	}
}

func configureExposedHeaders(options Options, w http.ResponseWriter) {
	if options.ExposeHeaders != "" {
		w.Header().Add(header.AccessControlExposeHeaders, options.ExposeHeaders)
	}
}

func configureCredentials(options Options, w http.ResponseWriter) {
	if options.AllowCredentials != "" {
		w.Header().Add(header.AccessControlAllowCredentials, options.AllowCredentials)
	}
}

func configureOrigin(options Options, r *http.Request, w http.ResponseWriter) {
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
