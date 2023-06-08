package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/httpd"
)

func TestHandlerReply(t *testing.T) {
	msg := "hello world"

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(msg))
	}

	m := httpd.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(httpd.Get(matcher.URLPath("/test")).
		Reply(httpd.Handler(fn)),
	)

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	req.Header.Add("test", "hello")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	txt, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
	assert.True(t, strings.Contains(string(txt), msg))
}
