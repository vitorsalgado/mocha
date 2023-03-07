package recover

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("content-type"))
	assert.True(t, strings.Contains(string(body), msg))
}
