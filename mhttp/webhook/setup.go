package webhook

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/mhttp"
)

// Builder provides a convenient way to configure a WebHook for a mock.
type Builder struct {
	input  *Input
	header http.Header
}

// Setup creates a new Builder to configure a WebHook for a mock.
func Setup() *Builder {
	return &Builder{input: &Input{Header: make(map[string]string)}, header: make(http.Header)}
}

// Method sets the HTTP method.
func (b *Builder) Method(method string) *Builder {
	b.input.Method = method
	return b
}

// URL sets the URL of the target service.
func (b *Builder) URL(u string) *Builder {
	b.input.URL = u
	return b
}

// Header adds an HTTP header to be sent in the WebHook.
func (b *Builder) Header(k, v string) *Builder {
	b.header.Add(k, v)
	return b
}

// SSLVerify enable or disable SSL verify for WebHook HTTP requests.
func (b *Builder) SSLVerify(v bool) *Builder {
	b.input.SSLVerify = v
	return b
}

// Body sets the HTTP request body of the WebHook.
func (b *Builder) Body(body string) *Builder {
	b.input.Body = body
	return b
}

// Transform sets a transformation function that allows users to customize the WebHook arguments.
func (b *Builder) Transform(transform Transform) *Builder {
	b.input.Transform = transform
	return b
}

// Build builds the mocha.PostActionDef required to set up a WebHook for a mock.
func (b *Builder) Build() *mhttp.PostActionDef {
	for k := range b.header {
		b.input.Header[k] = b.header.Get(k)
	}

	b.header = nil

	return &mhttp.PostActionDef{Name: Name, RawParameters: b.input}
}
