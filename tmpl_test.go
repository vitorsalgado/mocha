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

type testData struct {
	Key   string
	Value string
}

func TestGoTemplating(t *testing.T) {
	wd, _ := os.Getwd()
	filename := path.Join(wd, "testdata/test.tmpl")

	tpl, err := os.ReadFile(filename)
	require.NoError(t, err)

	tmpl := NewTextTemplate()
	err = tmpl.Template(string(tpl)).FuncMap(template.FuncMap{"trim": strings.TrimSpace}).Compile()
	require.NoError(t, err)

	data := testData{Key: "  hello   ", Value: "world "}
	buf := bytes.Buffer{}
	err = tmpl.Render(&buf, data)
	require.NoError(t, err)

	assert.Equal(t, "hello world \n", buf.String())
}

func TestTemplatingError(t *testing.T) {
	tmpl := NewTextTemplate()
	err := tmpl.Name("fail").Template("invalid {{ .hi }").Compile()

	assert.NotNil(t, err)
}

func TestReplyWithTemplate(t *testing.T) {
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

	res, err := NewReply().
		Status(http.StatusOK).
		BodyTemplate(NewTextTemplate().
			FuncMap(template.FuncMap{"trim": strings.TrimSpace}).
			Template(string(b)), data).
		Build(nil, newReqValues(req))

	require.NoError(t, err)

	b, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "test\ndev\n", string(b))
}
