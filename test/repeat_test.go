package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/v2"
	"github.com/vitorsalgado/mocha/v2/expect"
	"github.com/vitorsalgado/mocha/v2/internal/testutil"
	"github.com/vitorsalgado/mocha/v2/reply"
)

func TestRepeat(t *testing.T) {
	m := mocha.New(t)
	m.Start()

	m.AddMocks(mocha.Get(expect.URLPath("/test")).
		Repeat(3).
		Reply(reply.OK()))

	res, _ := testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusTeapot, res.StatusCode)
}
