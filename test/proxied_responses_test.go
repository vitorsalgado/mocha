package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestForward(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "ok", r.Header.Get("x-test"))
		assert.Equal(t, "", r.Header.Get("x-del"))
		assert.Equal(t, mimetype.TextPlain, r.Header.Get(header.ContentType))

		b, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			assert.NoError(t, err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))

	defer dest.Close()

	m := mocha.New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(mocha.Post(matcher.URLPath("/test")).
		Body(matcher.Equal("hello world")).
		Reply(mocha.From(dest.URL).
			ProxyHeader("x-test", "ok").
			Header("x-res", "example").
			RemoveProxyHeaders("x-del")))

	data := strings.NewReader("hello world")
	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", data)
	req.Header.Add("x-del", "to-delete")
	req.Header.Add(header.ContentType, mimetype.TextPlain)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.NoError(t, res.Body.Close())
	scoped.AssertCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "example", res.Header.Get("x-res"))
	assert.Equal(t, "hello world", string(b))
}
