package reply

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
		w := httptest.NewRecorder()
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

		res, err := Forward(dest.URL).
			ProxyHeader("x-proxy", "proxied").
			RemoveProxyHeader("x-to-be-removed").
			Header("x-res", "response").
			Build(w, newReqValues(req))

		assert.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "hello world", w.Body.String())
		assert.Equal(t, "response", w.Header().Get("x-res"))
	})

	t.Run("should forward and respond POST with body", func(t *testing.T) {
		w := httptest.NewRecorder()
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
		forward := Forward(u)
		res, err := forward.Build(w, newReqValues(req))

		assert.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("should forward and respond a No Content", func(t *testing.T) {
		w := httptest.NewRecorder()
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

		forward := Forward(dest.URL)
		res, err := forward.Build(w, newReqValues(req))

		assert.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("should remove prefix from URL", func(t *testing.T) {
		w := httptest.NewRecorder()
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/example", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))

			w.WriteHeader(http.StatusOK)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)
		res, err := Forward(dest.URL).TrimPrefix("/path/test").Build(w, newReqValues(req))

		assert.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should remove suffix from URL", func(t *testing.T) {
		w := httptest.NewRecorder()
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/path/test", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))

			w.WriteHeader(http.StatusOK)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)
		res, err := Forward(dest.URL).TrimSuffix("/example").Build(w, newReqValues(req))

		assert.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should panic if the given raw target cannot be parsed to a URL", func(t *testing.T) {
		assert.Panics(t, func() {
			Forward(" https://fail test  ")
		})
	})
}
