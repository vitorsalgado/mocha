package templating

import (
	"io/ioutil"
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

func TestGoTemplating_Compile(t *testing.T) {
	wd, _ := os.Getwd()
	filename := path.Join(wd, "_testdata/test.tmpl")

	tpl, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	tmpl := New()
	err = tmpl.Template(string(tpl)).FuncMap(template.FuncMap{"trim": strings.TrimSpace}).Compile()
	if err != nil {
		t.Fatal(err)
	}

	data := testData{Key: "  hello   ", Value: "world "}
	b, err := tmpl.Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello world \n", string(b))
}
