package templating

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"text/template"
)

type testData struct {
	Key   string
	Value string
}

func TestGoTemplating_Compile(t *testing.T) {
	wd, _ := os.Getwd()
	filename := path.Join(wd, "_testdata/test.tmpl")

	tpl, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	tmpl := New()
	tmpl.Template(string(tpl)).FuncMap(template.FuncMap{"trim": strings.TrimSpace})

	buf := &bytes.Buffer{}
	data := testData{Key: "  hello   ", Value: "world "}
	err = tmpl.Parse(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello world \n", buf.String())
}
