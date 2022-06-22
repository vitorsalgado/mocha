package mocha

import (
	"net/http"
	"strconv"
	"strings"
)

type CORSOptions struct {
	AllowedOrigin     string
	AllowCredentials  string
	AllowedMethods    string
	AllowedHeaders    string
	ExposeHeaders     string
	MaxAge            int
	SuccessStatusCode int
}

func CORS(options CORSOptions) func(http.Handler) http.Handler {
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

				w.Header().Add(HeaderVary, HeaderAccessControlRequestHeaders)
				w.Header().Add(HeaderContentLength, "0")

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

func configureHeaders(options CORSOptions, w http.ResponseWriter, r *http.Request) {
	// when allowed headers aren't specified, use values from header access-control-request-headers
	if options.AllowedHeaders != "" {
		w.Header().Add(HeaderAccessControlAllowHeaders, options.AllowedHeaders)
	} else {

		headers := r.Header.Get(HeaderAccessControlRequestHeaders)
		if strings.TrimSpace(headers) != "" {
			w.Header().Add(HeaderAccessControlAllowHeaders, headers)
		}
	}
}

func configureMaxAge(options CORSOptions, w http.ResponseWriter) {
	if options.MaxAge > -1 {
		w.Header().Add(HeaderAccessControlMaxAge, strconv.Itoa(options.MaxAge))
	}
}

func configureMethods(options CORSOptions, w http.ResponseWriter) {
	if len(options.AllowedMethods) > 0 {
		w.Header().Add(HeaderAccessControlAllowMethods, options.AllowedMethods)
	}
}

func configureExposedHeaders(options CORSOptions, w http.ResponseWriter) {
	if options.ExposeHeaders != "" {
		w.Header().Add(HeaderAccessControlExposeHeaders, options.ExposeHeaders)
	}
}

func configureCredentials(options CORSOptions, w http.ResponseWriter) {
	if options.AllowCredentials != "" {
		w.Header().Add(HeaderAccessControlAllowCredentials, options.AllowCredentials)
	}
}

func configureOrigin(options CORSOptions, r *http.Request, w http.ResponseWriter) {
	if options.AllowedOrigin == "" {
		return
	}

	origins := strings.Split(options.AllowedOrigin, ",")
	size := len(origins)

	if size == 1 {
		w.Header().Add(HeaderAccessControlAllowOrigin, options.AllowedOrigin)
		w.Header().Add(HeaderVary, HeaderOrigin)
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
		w.Header().Add(HeaderAccessControlAllowOrigin, origin)
		w.Header().Add(HeaderVary, HeaderOrigin)
	}
}

type CORSOptionsBuilder struct {
	options *CORSOptions
	origins []string
}

func CORSOpts() *CORSOptionsBuilder {
	return &CORSOptionsBuilder{
		origins: make([]string, 0),
		options: &CORSOptions{SuccessStatusCode: http.StatusNoContent}}
}

func (b *CORSOptionsBuilder) SuccessStatusCode(code int) *CORSOptionsBuilder {
	b.options.SuccessStatusCode = code
	return b
}

func (b *CORSOptionsBuilder) MaxAge(maxage int) *CORSOptionsBuilder {
	b.options.MaxAge = maxage
	return b
}

func (b *CORSOptionsBuilder) AllowOrigin(origin ...string) *CORSOptionsBuilder {
	b.origins = append(b.origins, origin...)
	return b
}

func (b *CORSOptionsBuilder) AllowCredentials(allow bool) *CORSOptionsBuilder {
	b.options.AllowCredentials = strconv.FormatBool(allow)
	return b
}

func (b *CORSOptionsBuilder) ExposeHeaders(headers ...string) *CORSOptionsBuilder {
	b.options.ExposeHeaders = strings.Join(headers, ",")
	return b
}

func (b *CORSOptionsBuilder) AllowedHeaders(headers ...string) *CORSOptionsBuilder {
	b.options.AllowedHeaders = strings.Join(headers, ",")
	return b
}

func (b *CORSOptionsBuilder) AllowMethods(methods ...string) *CORSOptionsBuilder {
	b.options.AllowedMethods = strings.Join(methods, ",")
	return b
}

func (b *CORSOptionsBuilder) Build() *CORSOptions {
	if len(b.origins) > 0 {
		if len(b.origins) == 1 {
			b.options.AllowedOrigin = b.origins[0]
		} else {
			b.options.AllowedOrigin = strings.Join(b.origins, ",")
		}
	}

	return b.options
}
