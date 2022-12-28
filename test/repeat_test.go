package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestRepeat(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Times(3).
		Reply(mocha.OK()))

	res, _ := testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	assert.Equal(t, mocha.StatusRequestDidNotMatch, res.StatusCode)
}
