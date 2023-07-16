package cors

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
	"github.com/vitorsalgado/mocha/v3/dzhttp/internal/mid"
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
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusNoContent, res.StatusCode)
		require.Equal(t, "*", res.Header.Get(httpval.HeaderAccessControlAllowOrigin))
		require.Equal(t, "x-expose-this", res.Header.Get(httpval.HeaderAccessControlExposeHeaders))
		require.Equal(t, "true", res.Header.Get(httpval.HeaderAccessControlAllowCredentials))
		require.Equal(t, "GET,POST", res.Header.Get(httpval.HeaderAccessControlAllowMethods))
		require.Equal(t, "x-allow-this,x-allow-that", res.Header.Get(httpval.HeaderAccessControlAllowHeaders))

		// check the actual request
		res, err = http.Get(ts.URL)
		require.NoError(t, err)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.True(t, strings.Contains(string(body), msg))
		require.Equal(t, "*", res.Header.Get(httpval.HeaderAccessControlAllowOrigin))
		require.Equal(t, "x-expose-this", res.Header.Get(httpval.HeaderAccessControlExposeHeaders))
		require.Equal(t, "true", res.Header.Get(httpval.HeaderAccessControlAllowCredentials))
		require.Equal(t, "text/plain", res.Header.Get("content-type"))
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
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
		require.Equal(t, "*", res.Header.Get(httpval.HeaderAccessControlAllowOrigin))
		require.Equal(t, "GET,POST", res.Header.Get(httpval.HeaderAccessControlAllowMethods))
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
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusNoContent, res.StatusCode)
		require.Equal(t, "", res.Header.Get(httpval.HeaderAccessControlAllowOrigin))
		require.Equal(t, "GET,POST", res.Header.Get(httpval.HeaderAccessControlAllowMethods))
	})

	t.Run("should not consider empty origin", func(t *testing.T) {
		ts := httptest.NewServer(
			mid.Compose(New(CORS().
				AllowOrigin("").
				build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusNoContent, res.StatusCode)
		require.Equal(t, "", res.Header.Get(httpval.HeaderAccessControlAllowOrigin))
	})
}
