package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

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
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}

func TestRepeat_Once(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Once().
		Reply(mocha.OK()))

	res, _ := testutil.Get(m.URL() + "/test").Do()
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = testutil.Get(m.URL() + "/test").Do()
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}
