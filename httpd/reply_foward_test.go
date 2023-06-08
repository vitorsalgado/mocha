package httpd

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/httpd/httpval"
)

func TestForward(t *testing.T) {
	t.Run("should forward and respond basic GET", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/path/test/example", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))
			assert.Equal(t, "", r.Header.Get("x-to-be-removed"))
			assert.Equal(t, "ok", r.Header.Get("x-present"))
			assert.Equal(t, []string{"proxied", "ok"}, r.Header.Values("x-proxy"))

			w.Header().Add("Trailer", "x-test-trailer")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("hello world"))
			w.Header().Add("x-test-trailer", "trailer-ok")
		}))

		defer dest.Close()

		h := make(http.Header)
		h.Add("x-res", "ok")

		ph := make(http.Header)
		ph.Add("x-proxy", "ok")

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)
		req.Header.Set("x-to-be-removed", "nok")
		req.Header.Set("x-present", "ok")
		rv := &RequestValues{RawRequest: req, URL: req.URL}

		reply := From(dest.URL).
			ForwardHeader("x-proxy", "proxied").
			ProxyHeaders(ph).
			RemoveProxyHeaders("x-to-be-removed").
			Header("x-res", "response").
			Headers(h)
		require.NoError(t, reply.beforeBuild(NewAPI()))
		res, err := reply.Build(nil, rv)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "hello world", string(res.Body))
		assert.Equal(t, []string{"response", "ok"}, res.Header.Values("x-res"))
		assert.Equal(t, "trailer-ok", res.Trailer.Get("x-test-trailer"))
	})

	t.Run("should forward and respond POST with body", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			b, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			w.Write(b)
		}))

		defer dest.Close()

		expected := "test text"
		body := strings.NewReader(expected)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", body)
		rv := &RequestValues{RawRequest: req, URL: req.URL}

		u, _ := url.Parse(dest.URL)
		forward := From(u)
		require.NoError(t, forward.beforeBuild(NewAPI()))
		res, err := forward.Build(nil, rv)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, expected, string(res.Body))
	})

	t.Run("should forward and respond POST with compressed body", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Add(httpval.HeaderContentEncoding, "gzip")
			w.Header().Add(httpval.HeaderContentType, "application/json")

			b, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			ww := gzip.NewWriter(w)

			defer ww.Close()

			_, _ = ww.Write(b)
		}))

		defer dest.Close()

		expected := "test text"
		body := strings.NewReader(expected)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", body)
		rv := &RequestValues{RawRequest: req, URL: req.URL}

		u, _ := url.Parse(dest.URL)
		forward := From(u)
		require.NoError(t, forward.beforeBuild(NewAPI()))
		res, err := forward.Build(nil, rv)
		require.NoError(t, err)

		g, err := res.Gunzip()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, expected, string(g))
	})

	t.Run("should forward and respond a No Content", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rv := &RequestValues{RawRequest: req, URL: req.URL}

		forward := From(dest.URL)
		require.NoError(t, forward.beforeBuild(NewAPI()))
		res, err := forward.Build(nil, rv)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "", string(res.Body))
	})

	t.Run("should remove prefix from URL", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/example", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))

			w.WriteHeader(http.StatusOK)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)
		rv := &RequestValues{RawRequest: req, URL: req.URL}
		reply := From(dest.URL).TrimPrefix("/path/test")
		require.NoError(t, reply.beforeBuild(NewAPI()))
		res, err := reply.Build(nil, rv)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should remove suffix from URL", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/path/test", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))

			w.WriteHeader(http.StatusOK)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)
		rv := &RequestValues{RawRequest: req, URL: req.URL}
		reply := From(dest.URL).TrimSuffix("/example")
		require.NoError(t, reply.beforeBuild(NewAPI()))
		res, err := reply.Build(nil, rv)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should panic if the given raw target cannot be parsed to a URL", func(t *testing.T) {
		assert.Panics(t, func() {
			From(" https://fail test  ")
		})
	})

	t.Run("init From with string and *url.URL", func(t *testing.T) {
		addr := "https://localhost:8080"
		u, _ := url.Parse(addr)

		assert.Equal(t, From(u).target.String(), addr)
		assert.Equal(t, From(addr).target.String(), addr)
	})

	t.Run("should forward and respond basic GET", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/path/test/example", r.URL.Path)

			<-time.After(500 * time.Millisecond)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("hello world"))
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example", nil)
		rv := &RequestValues{RawRequest: req, URL: req.URL}
		reply := From(dest.URL).
			Timeout(100 * time.Millisecond)
		require.NoError(t, reply.beforeBuild(NewAPI()))
		res, err := reply.Build(nil, rv)

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should follow redirects when using default configs", func(t *testing.T) {
		target := NewAPI()
		target.MustStart()

		defer target.Close()

		m := NewAPI()
		m.MustStart()

		defer m.Close()

		redirectScope := target.MustMock(Getf("/redirected").Reply(Accepted()))
		proxyScope := target.MustMock(Getf("/proxy").Reply(MovedPermanently(target.URL() + "/redirected")))
		scoped := m.MustMock(Getf("/proxy").Reply(From(target.URL())))

		res, err := http.Get(m.URL() + "/proxy")

		require.NoError(t, err)
		require.Equal(t, http.StatusAccepted, res.StatusCode)

		scoped.AssertCalled(t)
		proxyScope.AssertCalled(t)
		redirectScope.AssertCalled(t)
	})

	t.Run("should NOT follow redirects when it is disabled", func(t *testing.T) {
		target := NewAPI()
		target.MustStart()

		defer target.Close()

		m := NewAPI()
		m.MustStart()

		defer m.Close()

		redirectScope := target.MustMock(Getf("/redirected").Reply(Accepted()))
		proxyScope := target.MustMock(Getf("/proxy").Reply(MovedPermanently(target.URL() + "/redirected")))
		scoped := m.MustMock(Getf("/proxy").Reply(From(target.URL()).NoFollow()))

		httpClient := &http.Client{CheckRedirect: noFollow}
		res, err := httpClient.Get(m.URL() + "/proxy")

		require.NoError(t, err)
		require.Equal(t, http.StatusMovedPermanently, res.StatusCode)

		scoped.AssertCalled(t)
		proxyScope.AssertCalled(t)
		redirectScope.AssertNotCalled(t)
	})

	t.Run("tls", func(t *testing.T) {
		target := NewAPIWithT(t)
		target.MustStartTLS()

		m := NewAPIWithT(t)
		m.MustStart()

		s1 := target.MustMock(Getf("/test").Reply(OK().PlainText("hi")))
		s2 := m.MustMock(Getf("/test").Reply(From(target.URL()).SSLVerify(true)))

		httpClient := &http.Client{}
		res, err := httpClient.Get(m.URL() + "/test")
		require.NoError(t, err)

		defer res.Body.Close()

		s1.AssertNotCalled(t)
		s2.AssertNotCalled(t)

		require.Equal(t, StatusNoMatch, res.StatusCode)

		s2.Clean()
		s2 = m.MustMock(Getf("/test").Reply(From(target.URL()).SkipSSLVerify()))

		res, err = httpClient.Get(m.URL() + "/test")
		require.NoError(t, err)

		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		s1.AssertCalled(t)
		s2.AssertCalled(t)

		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, "hi", string(b))
	})
}
