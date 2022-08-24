package reply

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForward(t *testing.T) {
	t.Parallel()

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
			RemoveProxyHeader("x-to-be-removed").
			Header("x-res", "response").
			Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)
		assert.Equal(t, "hello world", string(b))
		assert.Equal(t, "response", res.Header.Get("x-res"))
	})

	t.Run("should forward and respond POST with body", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			b, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			w.Write(b)
		}))

		defer dest.Close()

		expected := "test text"
		body := strings.NewReader(expected)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", body)

		u, _ := url.Parse(dest.URL)
		forward := ProxiedFrom(u)
		res, err := forward.Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusOK, res.Status)
		assert.Equal(t, expected, string(b))
	})

	t.Run("should forward and respond a No Content", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

		forward := From(dest.URL)
		res, err := forward.Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusNoContent, res.Status)
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

		res, err := From(dest.URL).StripPrefix("/path/test").Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusOK, res.Status)
	})

	t.Run("should remove suffix from URL", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/path/test", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))

			w.WriteHeader(http.StatusOK)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)

		res, err := From(dest.URL).StripSuffix("/example").Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusOK, res.Status)
	})

	t.Run("should panic if provide raw target cannot be parsed to a URL", func(t *testing.T) {
		assert.Panics(t, func() {
			From(" http://fail test  ")
		})
	})
}
