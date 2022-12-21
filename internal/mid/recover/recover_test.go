package recover

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/x/event"
)

func TestRecover(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msg := "error test"
	fn := func(w http.ResponseWriter, r *http.Request) {
		panic(msg)
	}

	evt := event.New()
	evt.StartListening(ctx)

	ts := httptest.NewServer(New(t).Recover(http.HandlerFunc(fn)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusTeapot, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("content-type"))
	assert.True(t, strings.Contains(string(body), msg))
}
