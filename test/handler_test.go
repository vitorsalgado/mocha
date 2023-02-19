package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestHandlerReply(t *testing.T) {
	msg := "hello world"

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(msg))
	}

	m := mocha.New()
	m.MustStart()

	defer m.Close()

	m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Reply(mocha.Handler(fn)),
	)

	req := testutil.Get(m.URL() + "/test")
	req.Header("test", "hello")
	res, err := req.Do()
	require.NoError(t, err)

	txt, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
	assert.True(t, strings.Contains(string(txt), msg))
}
