package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestHandlerReply(t *testing.T) {
	key := "msg"
	msg := "hello world"

	fn := func(w http.ResponseWriter, r *http.Request) {
		arg := r.Context().Value(reply.KArg).(*reply.Arg)
		message, _ := arg.Params.Get(key)

		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(message.(string)))
	}

	m := mocha.New(t)
	m.Parameters().Set(key, msg)

	m.Start()

	defer m.Close()

	m.AddMocks(mocha.Get(matcher.URLPath("/test")).
		Reply(reply.Handler(fn)),
	)

	req := testutil.Get(m.URL() + "/test")
	req.Header("test", "hello")
	res, err := req.Do()
	assert.NoError(t, err)

	txt, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
	assert.True(t, strings.Contains(string(txt), msg))
}
