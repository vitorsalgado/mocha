package test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestFormUrlEncoded(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(mocha.Post(matcher.URLPath("/test")).
		FormField("var1", matcher.Equal("dev")).
		FormField("var2", matcher.Contain("q")).
		Reply(mocha.OK()))

	data := url.Values{}
	data.Set("var1", "dev")
	data.Set("var2", "qa")

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader(data.Encode()))
	req.Header.Add("test", "hello")
	req.Header.Add(header.ContentType, mimetype.FormURLEncoded)
	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.NoError(t, res.Body.Close())
	assert.True(t, scoped.HasBeenCalled())
}
