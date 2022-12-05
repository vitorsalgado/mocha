package reply

import (
	"io"
	"net/http"
	"text/template"
)

// Template defines a template parser for response bodies.
type Template interface {
	// Compile allows pre-compilation of the given template.
	Compile() error

	// Parse parses the given template.
	Parse(io.Writer, any) error
}

// TemplateData is the data model used to render the templates.
type TemplateData struct {
	// Request is HTTP request ref.
	Request *http.Request

	// Data is the model to be used with the given template.
	// This value is set using the Model() function from StdReply.
	Data any
}

// TextTemplate is the built-in text Template interface.
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

func (gt *TextTemplate) Parse(w io.Writer, data any) error {
	return gt.t.Execute(w, data)
}
