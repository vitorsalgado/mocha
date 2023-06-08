package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/httpd"
)

func TestRepeat(t *testing.T) {
	m := mhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(mhttp.Get(matcher.URLPath("/test")).
		Times(3).
		Reply(mhttp.OK()))

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
	require.Equal(t, mhttp.StatusNoMatch, res.StatusCode)
	require.NoError(t, err)
}

func TestRepeat_Once(t *testing.T) {
	m := mhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(mhttp.Get(matcher.URLPath("/test")).
		Once().
		Reply(mhttp.OK()))

	res, _ := http.Get(m.URL() + "/test")
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = http.Get(m.URL() + "/test")
	require.Equal(t, mhttp.StatusNoMatch, res.StatusCode)
}

func TestRepeat_FileSetup(t *testing.T) {
	m := mhttp.NewAPIWithT(t)
	m.MustStart()
	m.MustMock(mhttp.FromFile("testdata/repeat/1_repeat.yaml"))

	httpClient := &http.Client{}

	for i := 0; i < 5; i++ {
		res, err := httpClient.Get(m.URL("test"))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode, i)
	}

	res, err := httpClient.Get(m.URL("/test"))
	require.NoError(t, err)
	require.Equal(t, mhttp.StatusNoMatch, res.StatusCode)
}
