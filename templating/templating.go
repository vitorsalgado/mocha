package templating

import (
	"io"
	"text/template"
)

type (
	Template interface {
		Compile() error
		Parse(io.Writer, any) error
	}

	BuiltInTemplate struct {
		name     string
		funcMap  template.FuncMap
		template string
		t        *template.Template
	}
)

func New() *BuiltInTemplate {
	return &BuiltInTemplate{funcMap: make(template.FuncMap)}
}

func (gt *BuiltInTemplate) Name(name string) *BuiltInTemplate {
	gt.name = name
	return gt
}

func (gt *BuiltInTemplate) FuncMap(fn template.FuncMap) *BuiltInTemplate {
	gt.funcMap = fn
	return gt
}

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
