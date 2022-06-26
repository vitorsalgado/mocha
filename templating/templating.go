package templating

import (
	"io"
	"text/template"
)

type (
	Parser interface {
		Parse(io.Writer, any) error
	}

	GoTemplateParser struct {
		name     string
		funcMap  template.FuncMap
		template string
	}
)

func New() *GoTemplateParser {
	return &GoTemplateParser{funcMap: make(template.FuncMap)}
}

func (gt *GoTemplateParser) Name(name string) *GoTemplateParser {
	gt.name = name
	return gt
}

func (gt *GoTemplateParser) FuncMap(fn template.FuncMap) *GoTemplateParser {
	gt.funcMap = fn
	return gt
}

func (gt *GoTemplateParser) Template(tmpl string) *GoTemplateParser {
	gt.template = tmpl
	return gt
}

func (gt *GoTemplateParser) Parse(w io.Writer, data any) error {
	t, err := template.New(gt.name).Funcs(gt.funcMap).Parse(gt.template)
	if err != nil {
		return err
	}

	return t.Execute(w, data)
}
