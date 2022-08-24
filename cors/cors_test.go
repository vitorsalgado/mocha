package cors

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v2/internal/headers"
	"github.com/vitorsalgado/mocha/v2/internal/middleware"
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
			middleware.Compose(New(Configure().
				AllowMethods("GET", "POST").
				AllowedHeaders("x-allow-this", "x-allow-that").
				ExposeHeaders("x-expose-this").
				AllowOrigin("*").
				MaxAge(10).
				AllowCredentials(true).
				Build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		// check preflight request
		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "*", res.Header.Get(headers.AccessControlAllowOrigin))
		assert.Equal(t, "x-expose-this", res.Header.Get(headers.AccessControlExposeHeaders))
		assert.Equal(t, "true", res.Header.Get(headers.AccessControlAllowCredentials))
		assert.Equal(t, "GET,POST", res.Header.Get(headers.AccessControlAllowMethods))
		assert.Equal(t, "x-allow-this,x-allow-that", res.Header.Get(headers.AccessControlAllowHeaders))

		// check the actual request
		res, err = http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, strings.Contains(string(body), msg))
		assert.Equal(t, "*", res.Header.Get(headers.AccessControlAllowOrigin))
		assert.Equal(t, "x-expose-this", res.Header.Get(headers.AccessControlExposeHeaders))
		assert.Equal(t, "true", res.Header.Get(headers.AccessControlAllowCredentials))
		assert.Equal(t, "text/plain", res.Header.Get("content-type"))
	})

	t.Run("should return custom success status code", func(t *testing.T) {
		ts := httptest.NewServer(
			middleware.Compose(New(Configure().
				AllowMethods("GET", "POST").
				AllowOrigin("*").
				SuccessStatusCode(http.StatusBadRequest).
				Build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "*", res.Header.Get(headers.AccessControlAllowOrigin))
		assert.Equal(t, "GET,POST", res.Header.Get(headers.AccessControlAllowMethods))
	})

	t.Run("should check origin from a list when one is provided", func(t *testing.T) {
		ts := httptest.NewServer(
			middleware.Compose(New(Configure().
				AllowMethods("GET", "POST").
				AllowOrigin("http://localhost:8080", "http://localhost:8081").
				Build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "", res.Header.Get(headers.AccessControlAllowOrigin))
		assert.Equal(t, "GET,POST", res.Header.Get(headers.AccessControlAllowMethods))
	})

	t.Run("should not consider empty origin", func(t *testing.T) {
		ts := httptest.NewServer(
			middleware.Compose(New(Configure().
				AllowOrigin("").
				Build())).Root(http.HandlerFunc(handler)))
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "", res.Header.Get(headers.AccessControlAllowOrigin))
	})
}
