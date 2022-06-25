package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type OptionsBuilder struct {
	options *Options
	origins []string
}

func Configure() *OptionsBuilder {
	return &OptionsBuilder{
		origins: make([]string, 0),
		options: &Options{SuccessStatusCode: http.StatusNoContent}}
}

func (b *OptionsBuilder) SuccessStatusCode(code int) *OptionsBuilder {
	b.options.SuccessStatusCode = code
	return b
}

func (b *OptionsBuilder) MaxAge(maxAge int) *OptionsBuilder {
	b.options.MaxAge = maxAge
	return b
}

func (b *OptionsBuilder) AllowOrigin(origin ...string) *OptionsBuilder {
	b.origins = append(b.origins, origin...)
	return b
}

func (b *OptionsBuilder) AllowCredentials(allow bool) *OptionsBuilder {
	b.options.AllowCredentials = strconv.FormatBool(allow)
	return b
}

func (b *OptionsBuilder) ExposeHeaders(headers ...string) *OptionsBuilder {
	b.options.ExposeHeaders = strings.Join(headers, ",")
	return b
}

func (b *OptionsBuilder) AllowedHeaders(headers ...string) *OptionsBuilder {
	b.options.AllowedHeaders = strings.Join(headers, ",")
	return b
}

func (b *OptionsBuilder) AllowMethods(methods ...string) *OptionsBuilder {
	b.options.AllowedMethods = strings.Join(methods, ",")
	return b
}

func (b *OptionsBuilder) Build() *Options {
	if len(b.origins) > 0 {
		if len(b.origins) == 1 {
			b.options.AllowedOrigin = b.origins[0]
		} else {
			b.options.AllowedOrigin = strings.Join(b.origins, ",")
		}
	}

	return b.options
}
