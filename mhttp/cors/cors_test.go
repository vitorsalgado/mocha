package cors

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/mhttp/internal/mid"
	"github.com/vitorsalgado/mocha/v3/misc"
)

func TestCORS(t *testing.T) {
	msg := "hello world"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	}

	t.Run("should allow request", func(t *testing.T) {
		ts := httptest.NewServer(
			mid.Compose(New(CORS().
				AllowMethods("GET", "POST").
				AllowedHeaders("x-allow-this", "x-allow-that").
				ExposeHeaders("x-expose-this").
				AllowOrigin("*").
				MaxAge(10).
				AllowCredentials(true).
				build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		// check preflight request
		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "*", res.Header.Get(misc.HeaderAccessControlAllowOrigin))
		assert.Equal(t, "x-expose-this", res.Header.Get(misc.HeaderAccessControlExposeHeaders))
		assert.Equal(t, "true", res.Header.Get(misc.HeaderAccessControlAllowCredentials))
		assert.Equal(t, "GET,POST", res.Header.Get(misc.HeaderAccessControlAllowMethods))
		assert.Equal(t, "x-allow-this,x-allow-that", res.Header.Get(misc.HeaderAccessControlAllowHeaders))

		// check the actual request
		res, err = http.Get(ts.URL)
		assert.NoError(t, err)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, strings.Contains(string(body), msg))
		assert.Equal(t, "*", res.Header.Get(misc.HeaderAccessControlAllowOrigin))
		assert.Equal(t, "x-expose-this", res.Header.Get(misc.HeaderAccessControlExposeHeaders))
		assert.Equal(t, "true", res.Header.Get(misc.HeaderAccessControlAllowCredentials))
		assert.Equal(t, "text/plain", res.Header.Get("content-type"))
	})

	t.Run("should return custom success status code", func(t *testing.T) {
		ts := httptest.NewServer(
			mid.Compose(New(CORS().
				AllowMethods("GET", "POST").
				AllowOrigin("*").
				SuccessStatusCode(http.StatusBadRequest).
				build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "*", res.Header.Get(misc.HeaderAccessControlAllowOrigin))
		assert.Equal(t, "GET,POST", res.Header.Get(misc.HeaderAccessControlAllowMethods))
	})

	t.Run("should check origin from a list when one is provided", func(t *testing.T) {
		ts := httptest.NewServer(
			mid.Compose(New(CORS().
				AllowMethods("GET", "POST").
				AllowOrigin("http://localhost:8080", "http://localhost:8081").
				build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "", res.Header.Get(misc.HeaderAccessControlAllowOrigin))
		assert.Equal(t, "GET,POST", res.Header.Get(misc.HeaderAccessControlAllowMethods))
	})

	t.Run("should not consider empty origin", func(t *testing.T) {
		ts := httptest.NewServer(
			mid.Compose(New(CORS().
				AllowOrigin("").
				build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "", res.Header.Get(misc.HeaderAccessControlAllowOrigin))
	})
}
