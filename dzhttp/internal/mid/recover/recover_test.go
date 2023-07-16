package recover

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	msg := "error test"
	fn := func(w http.ResponseWriter, r *http.Request) {
		panic(msg)
	}

	ts := httptest.NewServer(New(func(_ error) {

	}, http.StatusInternalServerError).Recover(http.HandlerFunc(fn)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, "text/plain; charset=utf-8", res.Header.Get("content-type"))
	require.True(t, strings.Contains(string(body), msg))
}
