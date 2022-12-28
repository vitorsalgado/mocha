package mocha

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForward(t *testing.T) {
	t.Run("should forward and respond basic GET", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/path/test/example", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))
			assert.Equal(t, "", r.Header.Get("x-to-be-removed"))
			assert.Equal(t, "ok", r.Header.Get("x-present"))
			assert.Equal(t, "proxied", r.Header.Get("x-proxy"))

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("hello world"))
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)
		req.Header.Set("x-to-be-removed", "nok")
		req.Header.Set("x-present", "ok")

		res, err := From(dest.URL).
			ProxyHeader("x-proxy", "proxied").
			RemoveProxyHeaders("x-to-be-removed").
			Header("x-res", "response").
			Build(nil, newReqValues(req))

		require.NoError(t, err)

		b, err := io.ReadAll(res.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "hello world", string(b))
		assert.Equal(t, "response", res.Header.Get("x-res"))
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

		u, _ := url.Parse(dest.URL)
		forward := From(u)
		res, err := forward.Build(nil, newReqValues(req))

		require.NoError(t, err)

		b, err := io.ReadAll(res.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, expected, string(b))
	})

	t.Run("should forward and respond a No Content", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

		forward := From(dest.URL)
		res, err := forward.Build(nil, newReqValues(req))

		require.NoError(t, err)

		b, err := io.ReadAll(res.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "", string(b))
	})

	t.Run("should remove prefix from URL", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/example", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))

			w.WriteHeader(http.StatusOK)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)
		res, err := From(dest.URL).TrimPrefix("/path/test").Build(nil, newReqValues(req))

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
		res, err := From(dest.URL).TrimSuffix("/example").Build(nil, newReqValues(req))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should panic if the given raw target cannot be parsed to a URL", func(t *testing.T) {
		assert.Panics(t, func() {
			From(" https://fail test  ")
		})
	})

	t.Run("should accept url.URL pointer and non-pointer", func(t *testing.T) {
		addr := "https://localhost:8080"
		u, _ := url.Parse(addr)

		assert.Equal(t, From(u).target.String(), addr)
		assert.Equal(t, From(*u).target.String(), addr)
		assert.Equal(t, From(addr).target.String(), addr)
	})
}
