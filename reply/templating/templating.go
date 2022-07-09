package templating

import (
	"io"
	"net/http"
	"text/template"
)

type (
	// Template defines a template parser for response bodies.
	Template interface {
		// Compile allows pre-compilation of the given template.
		Compile() error

		// Parse parses the given template.
		Parse(io.Writer, any) error
	}

	Model struct {
		Request *http.Request
		Data    any
	}

	// BuiltInTemplate is the built-in implementation for Template interface.
	// It uses Go templates.
	BuiltInTemplate struct {
		name     string
		funcMap  template.FuncMap
		template string
		t        *template.Template
	}
)

// New creates a new BuiltInTemplate.
func New() *BuiltInTemplate {
	return &BuiltInTemplate{funcMap: make(template.FuncMap)}
}

// Name sets the name of the template.
func (gt *BuiltInTemplate) Name(name string) *BuiltInTemplate {
	gt.name = name
	return gt
}

// FuncMap adds a new function to be used inside the Go template.
func (gt *BuiltInTemplate) FuncMap(fn template.FuncMap) *BuiltInTemplate {
	gt.funcMap = fn
	return gt
}

// Template sets the actual template.
func (gt *BuiltInTemplate) Template(tmpl string) *BuiltInTemplate {
	gt.template = tmpl
	return gt
}

func (gt *BuiltInTemplate) Compile() error {
	t, err := template.New(gt.name).Funcs(gt.funcMap).Parse(gt.template)
	if err != nil {
		return err
	}

	gt.t = t

	return nil
}

func (gt *BuiltInTemplate) Parse(w io.Writer, data any) error {
	return gt.t.Execute(w, data)
}
