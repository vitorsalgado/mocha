package mocha

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoTemplating(t *testing.T) {
	type testData struct {
		Key   string
		Value string
	}

	wd, _ := os.Getwd()
	filename := path.Join(wd, "testdata/test.tmpl")

	tpl, err := os.ReadFile(filename)
	require.NoError(t, err)

	gt := newGoTemplate()
	tr, err := gt.FuncMap(template.FuncMap{"trim": strings.TrimSpace}).Parse(string(tpl))
	require.NoError(t, err)

	data := testData{Key: "  hello   ", Value: "world "}
	buf := bytes.Buffer{}
	err = tr.Render(&buf, data)

	require.NoError(t, err)
	require.Equal(t, "hello world \n", buf.String())
}

func TestTemplatingError(t *testing.T) {
	gt := newGoTemplate()
	tr, err := gt.Parse("invalid {{ .hi }")

	require.NotNil(t, err)
	require.Nil(t, tr)
}

func TestReplyWithTemplate(t *testing.T) {
	app := New(Configure().TemplateEngineFunctions(template.FuncMap{"trim": strings.TrimSpace}))
	app.MustStart()
	defer app.Close()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	req.Header.Add("x-test", "dev")

	wd, _ := os.Getwd()
	f, _ := os.Open(path.Join(wd, "testdata/test_req.tmpl"))
	b, _ := io.ReadAll(f)

	data := struct {
		Name string
	}{
		Name: " test  ",
	}

	rv := &RequestValues{RawRequest: req, URL: req.URL}
	reply := NewReply().
		Status(http.StatusOK).
		BodyTemplate(string(b)).
		SetTemplateData(data)

	require.NoError(t, reply.beforeBuild(app))
	res, err := reply.Build(nil, rv)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "test\ndev\n", string(res.Body))
}

func TestReplyWithTemplateText(t *testing.T) {
	app := New(Configure().TemplateEngineFunctions(template.FuncMap{"trim": strings.TrimSpace}))
	app.MustStart()
	defer app.Close()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	req.Header.Add("x-test", "dev")

	data := struct {
		Name string
	}{
		Name: " test  ",
	}

	tmpl := `{{- trim .Ext.Name }}
{{ .Request.Header.Get "x-test" }}
`

	rv := &RequestValues{RawRequest: req, URL: req.URL}
	reply := NewReply().
		Status(http.StatusOK).
		BodyTemplate(tmpl).
		SetTemplateData(data)
	require.NoError(t, reply.beforeBuild(app))
	res, err := reply.Build(nil, rv)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "test\ndev\n", string(res.Body))
}
