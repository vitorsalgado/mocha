package mocha

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
)

var (
	_ TemplateEngine   = (*builtInGoTemplate)(nil)
	_ TemplateRenderer = (*builtInGoTemplateRender)(nil)
)

type TemplateEngine interface {
	// Load is executed on server initialization.
	Load() error

	// Parse pre-compiles the source template. It will be called once during mock setup.
	Parse(string) (TemplateRenderer, error)
}

// TemplateRenderer defines a template parser for response bodies.
type TemplateRenderer interface {
	// Render renders the previously parsed template to the given io.Writer.
	// The second parameter is template data.
	Render(io.Writer, any) error
}

// templateData is the data tmplExt used to render the templates.
type templateData struct {
	// Request is HTTP request ref.
	Request templateRequest

	// App is the server instance.
	App *Mocha

	// Ext is an additional data that can be passed to the template.
	Ext any
}

type templateRequest struct {
	Method          string
	URL             *url.URL
	URLPathSegments []string
	Header          http.Header
	Cookies         []*http.Cookie
	Body            any
}

func (t *templateRequest) Cookie(name string) (*http.Cookie, error) {
	for _, cookie := range t.Cookies {
		if cookie.Name == name {
			return cookie, nil
		}
	}

	return nil, fmt.Errorf(`cookie "%s" not found`, name)
}

// builtInGoTemplate is the built-in Template implementation.
// It uses Go templates.
type builtInGoTemplate struct {
	funcMap template.FuncMap
}

// newGoTemplate creates a new BuiltInTemplate.
func newGoTemplate() *builtInGoTemplate {
	return &builtInGoTemplate{funcMap: make(template.FuncMap)}
}

func (gt *builtInGoTemplate) Load() error {
	return nil
}

// FuncMap adds a new function to be used inside the Go template.
func (gt *builtInGoTemplate) FuncMap(fn template.FuncMap) *builtInGoTemplate {
	gt.funcMap = fn
	return gt
}

func (gt *builtInGoTemplate) Parse(s string) (TemplateRenderer, error) {
	t, err := template.New("").Funcs(gt.funcMap).Parse(s)
	if err != nil {
		return nil, err
	}

	return newGoTemplateRender(t), nil
}

type builtInGoTemplateRender struct {
	template *template.Template
}

func newGoTemplateRender(template *template.Template) TemplateRenderer {
	return &builtInGoTemplateRender{template: template}
}

func (gt *builtInGoTemplateRender) Render(w io.Writer, data any) error {
	return gt.template.Execute(w, data)
}
