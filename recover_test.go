package mocha

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msg := "error test"
	fn := func(w http.ResponseWriter, r *http.Request) {
		panic(msg)
	}

	evt := newEvents()
	evt.StartListening(ctx)

	rm := &recoverMid{d: func(err error) {}, t: t, evt: evt}
	ts := httptest.NewServer(rm.Recover(http.HandlerFunc(fn)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusTeapot, res.StatusCode)
	assert.Equal(t, "text/plain", res.Header.Get("content-type"))
	assert.True(t, strings.Contains(string(body), msg))
}
