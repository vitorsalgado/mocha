package reply

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
)

type testData struct {
	Key   string
	Value string
}

func TestGoTemplating(t *testing.T) {
	wd, _ := os.Getwd()
	filename := path.Join(wd, "testdata/test.tmpl")

	tpl, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	tmpl := NewTextTemplate()
	err = tmpl.Template(string(tpl)).FuncMap(template.FuncMap{"trim": strings.TrimSpace}).Compile()
	if err != nil {
		t.Fatal(err)
	}

	data := testData{Key: "  hello   ", Value: "world "}
	buf := bytes.Buffer{}
	err = tmpl.Parse(&buf, data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello world \n", buf.String())
}

func TestTemplatingError(t *testing.T) {
	tmpl := NewTextTemplate()
	err := tmpl.Name("fail").Template("invalid {{ .hi }").Compile()

	assert.NotNil(t, err)
}

func TestReplyWithTemplate(t *testing.T) {
	_req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	_req.Header.Add("x-test", "dev")

	wd, _ := os.Getwd()
	f, _ := os.Open(path.Join(wd, "testdata/test_req.tmpl"))
	b, _ := io.ReadAll(f)

	data := struct {
		Name string
	}{
		Name: " test  ",
	}

	res, err := New().
		Status(http.StatusOK).
		BodyTemplate(NewTextTemplate().
			FuncMap(template.FuncMap{"trim": strings.TrimSpace}).
			Template(string(b))).
		BodyTemplateModel(data).
		Build(nil, _req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.Status)
	assert.Equal(t, "test\ndev\n", string(res.Body))
}
