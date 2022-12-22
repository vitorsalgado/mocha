package reply

import (
	"io"
	"net/http"
	"net/url"
	"text/template"
)

// templateData is the data templateExtras used to render the templates.
type templateData struct {
	// Request is HTTP request ref.
	Request templateRequest

	// Extras is an additional data that can be passed to the template.
	// This value is set using the TemplateExtra() function from StdReply.
	Extras any
}

type templateRequest struct {
	Method string
	URL    url.URL
	Header http.Header
	Body   any
}

// Template defines a template parser for response bodies.
type Template interface {
	// Compile allows pre-compilation of the given template.
	Compile() error

	// Render parses the given template.
	Render(io.Writer, any) error
}

// TextTemplate is the built-in Template implementation.
// It uses Go templates.
type TextTemplate struct {
	name     string
	funcMap  template.FuncMap
	template string
	t        *template.Template
}

// NewTextTemplate creates a new BuiltInTemplate.
func NewTextTemplate() *TextTemplate {
	return &TextTemplate{funcMap: make(template.FuncMap)}
}

// Name sets the name of the template.
func (gt *TextTemplate) Name(name string) *TextTemplate {
	gt.name = name
	return gt
}

// FuncMap adds a new function to be used inside the Go template.
func (gt *TextTemplate) FuncMap(fn template.FuncMap) *TextTemplate {
	gt.funcMap = fn
	return gt
}

// Template sets the actual template.
func (gt *TextTemplate) Template(tmpl string) *TextTemplate {
	gt.template = tmpl
	return gt
}

func (gt *TextTemplate) Compile() error {
	t, err := template.New(gt.name).Funcs(gt.funcMap).Parse(gt.template)
	if err != nil {
		return err
	}

	gt.t = t

	return nil
}

func (gt *TextTemplate) Render(w io.Writer, data any) error {
	return gt.t.Execute(w, data)
}
