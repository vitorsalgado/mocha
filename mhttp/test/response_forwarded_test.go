package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
	"github.com/vitorsalgado/mocha/v3/misc"
)

func TestProxiedReplies(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "ok", r.Header.Get("x-test"))
		assert.Equal(t, "", r.Header.Get("x-del"))
		assert.Equal(t, misc.MIMETextPlain, r.Header.Get(misc.HeaderContentType))

		b, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			require.NoError(t, err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))

	defer dest.Close()

	m := mhttp2.NewAPI()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(mhttp2.Post(URLPath("/test")).
		Body(StrictEqual("hello world")).
		Reply(mhttp2.From(dest.URL).
			ForwardHeader("x-test", "ok").
			Header("x-res", "example").
			RemoveProxyHeaders("x-del")))

	data := strings.NewReader("hello world")
	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", data)
	req.Header.Add("x-del", "to-delete")
	req.Header.Add(misc.HeaderContentType, misc.MIMETextPlain)

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
	target := mhttp2.NewAPI()
	target.MustStart()
	defer target.Close()

	targetScoped := target.MustMock(
		mhttp2.Postf("/test").
			Headerf("test", "ok").
			Header("del", Not(Present())).
			Reply(mhttp2.OK().PlainText("done")))

	m := mhttp2.NewAPI()
	m.MustStart()
	defer m.Close()

	data := make(map[string]any)
	data["target"] = target.URL()
	m.SetData(data)

	scoped := m.MustMock(mhttp2.FromFile("testdata/response_forwarded/proxied_response.yaml"))
	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader("hello world"))
	req.Header.Add(misc.HeaderContentType, misc.MIMETextPlain)
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