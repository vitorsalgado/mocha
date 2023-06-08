package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/httpd"
	"github.com/vitorsalgado/mocha/v3/httpd/httpval"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestProxiedReplies(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "ok", r.Header.Get("x-test"))
		assert.Equal(t, "", r.Header.Get("x-del"))
		assert.Equal(t, httpval.MIMETextPlain, r.Header.Get(httpval.HeaderContentType))

		b, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			require.NoError(t, err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))

	defer dest.Close()

	m := httpd.NewAPI()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(httpd.Post(URLPath("/test")).
		Body(StrictEqual("hello world")).
		Reply(httpd.From(dest.URL).
			ForwardHeader("x-test", "ok").
			Header("x-res", "example").
			RemoveProxyHeaders("x-del")))

	data := strings.NewReader("hello world")
	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", data)
	req.Header.Add("x-del", "to-delete")
	req.Header.Add(httpval.HeaderContentType, httpval.MIMETextPlain)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.True(t, scoped.AssertCalled(t))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "example", res.Header.Get("x-res"))
	assert.Equal(t, "hello world", string(b))
}

func TestProxiedReplyMockFileWithTemplate(t *testing.T) {
	target := httpd.NewAPI()
	target.MustStart()
	defer target.Close()

	targetScoped := target.MustMock(
		httpd.Postf("/test").
			Headerf("test", "ok").
			Header("del", Not(Present())).
			Reply(httpd.OK().PlainText("done")))

	m := httpd.NewAPI()
	m.MustStart()
	defer m.Close()

	data := make(map[string]any)
	data["target"] = target.URL()
	m.SetData(data)

	scoped := m.MustMock(httpd.FromFile("testdata/response_forwarded/proxied_response.yaml"))
	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader("hello world"))
	req.Header.Add(httpval.HeaderContentType, httpval.MIMETextPlain)
	req.Header.Add("del", "to be deleted")

	res, err := httpClient.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "done", string(b))
	require.True(t, targetScoped.AssertCalled(t))
	require.True(t, scoped.AssertCalled(t))
}
