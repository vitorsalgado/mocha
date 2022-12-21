package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestRepeat(t *testing.T) {
	m := mocha.NewWithT(t)
	m.MustStart()

	defer m.Close()

	m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Times(3).
		Reply(reply.OK()))

	res, _ := testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, mocha.StatusRequestDidNotMatch, res.StatusCode)
}
