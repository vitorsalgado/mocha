package test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestFormUrlEncoded(t *testing.T) {
	m := mocha.New(t)
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(mocha.Post(expect.URLPath("/test")).
		FormField("var1", expect.ToEqual("dev")).
		FormField("var2", expect.ToContain("q")).
		Reply(reply.OK()))

	data := url.Values{}
	data.Set("var1", "dev")
	data.Set("var2", "qa")

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader(data.Encode()))
	req.Header.Add("test", "hello")
	req.Header.Add(header.ContentType, mimetype.FormURLEncoded)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	defer res.Body.Close()

	assert.True(t, scoped.Called())
}
