package dzhttp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"

	"github.com/vitorsalgado/mocha/v3/dzstd"
)

var (
	_ dzstd.TemplateEngine   = (*builtInGoTemplate)(nil)
	_ dzstd.TemplateRenderer = (*builtInGoTemplateRender)(nil)
)

// templateData is the data used in templates during rendering.
type templateData struct {
	Request templateRequest
	App     *templateAppWrapper
	Ext     any
}

type templateAppWrapper struct {
	app *HTTPMockApp
}

func (t *templateAppWrapper) Parameters() dzstd.Params {
	return t.app.Parameters()
}

func (t *templateAppWrapper) Data() map[string]any {
	return t.app.Data()
}

func (t *templateAppWrapper) URL(paths ...string) string {
	return t.app.URL(paths...)
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

// builtInGoTemplate is the built-in TemplateEngine implementation.
// It uses Go templates.
type builtInGoTemplate struct {
	funcMap template.FuncMap
}

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

func (gt *builtInGoTemplate) Parse(s string) (dzstd.TemplateRenderer, error) {
	t, err := template.New("").Funcs(gt.funcMap).Parse(s)
	if err != nil {
		return nil, err
	}

	return newGoTemplateRender(t), nil
}

type builtInGoTemplateRender struct {
	template *template.Template
}

func newGoTemplateRender(template *template.Template) dzstd.TemplateRenderer {
	return &builtInGoTemplateRender{template: template}
}

func (gt *builtInGoTemplateRender) Render(w io.Writer, data any) error {
	return gt.template.Execute(w, data)
}
