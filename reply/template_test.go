package reply

import (
	"bytes"
	"io/ioutil"
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
	filename := path.Join(wd, "_testdata/test.tmpl")

	tpl, err := ioutil.ReadFile(filename)
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
	req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	req.Header.Add("x-test", "dev")

	wd, _ := os.Getwd()
	f, _ := os.Open(path.Join(wd, "_testdata/test_req.tmpl"))
	b, _ := ioutil.ReadAll(f)

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
		Model(data).
		Build(req, &testMock, nil)

	if err != nil {
		t.Fatal(err)
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.Status)
	assert.Equal(t, "test\ndev\n", string(b))
}
