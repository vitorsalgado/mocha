package mid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddlewaresComposition(t *testing.T) {
	msg := "hello world"
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("x-one", r.Header.Get("x-one"))
		w.Header().Add("x-two", r.Header.Get("x-two"))
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	}

	ts := httptest.NewServer(Compose(one, two).Root(http.HandlerFunc(fn)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	require.NoError(t, err)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "ok", res.Header.Get("x-one"))
	require.Equal(t, "nok", res.Header.Get("x-two"))
	require.Equal(t, "text/plain", res.Header.Get("content-type"))
	require.True(t, strings.Contains(string(body), msg))
}

func one(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("x-one", "ok")
		next.ServeHTTP(w, r)
	})
}

func two(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("x-two", "nok")
		next.ServeHTTP(w, r)
	})
}
