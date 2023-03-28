package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
)

func TestRepeat(t *testing.T) {
	m := mhttp2.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(mhttp2.Get(matcher.URLPath("/test")).
		Times(3).
		Reply(mhttp2.OK()))

	res, err := http.Get(m.URL() + "/test")
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, err)

	res, err = http.Get(m.URL() + "/test")
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, err)

	res, err = http.Get(m.URL() + "/test")
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, err)

	res, err = http.Get(m.URL() + "/test")
	require.Equal(t, mhttp2.StatusNoMatch, res.StatusCode)
	require.NoError(t, err)
}

func TestRepeat_Once(t *testing.T) {
	m := mhttp2.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(mhttp2.Get(matcher.URLPath("/test")).
		Once().
		Reply(mhttp2.OK()))

	res, _ := http.Get(m.URL() + "/test")
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = http.Get(m.URL() + "/test")
	require.Equal(t, mhttp2.StatusNoMatch, res.StatusCode)
}

func TestRepeat_FileSetup(t *testing.T) {
	m := mhttp2.NewAPIWithT(t)
	m.MustStart()
	m.MustMock(mhttp2.FromFile("testdata/repeat/1_repeat.yaml"))

	httpClient := &http.Client{}

	for i := 0; i < 5; i++ {
		res, err := httpClient.Get(m.URL("test"))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode, i)
	}

	res, err := httpClient.Get(m.URL("/test"))
	require.NoError(t, err)
	require.Equal(t, mhttp2.StatusNoMatch, res.StatusCode)
}
