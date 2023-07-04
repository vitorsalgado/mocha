package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestRepeat(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(dzhttp.Get(matcher.URLPath("/test")).
		Times(3).
		Reply(dzhttp.OK()))

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
	require.Equal(t, dzhttp.StatusNoMatch, res.StatusCode)
	require.NoError(t, err)
}

func TestRepeat_Once(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(dzhttp.Get(matcher.URLPath("/test")).
		Once().
		Reply(dzhttp.OK()))

	res, _ := http.Get(m.URL() + "/test")
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = http.Get(m.URL() + "/test")
	require.Equal(t, dzhttp.StatusNoMatch, res.StatusCode)
}

func TestRepeat_FileSetup(t *testing.T) {
	m := dzhttp.NewAPIWithT(t)
	m.MustStart()
	m.MustMock(dzhttp.FromFile("testdata/repeat/1_repeat.yaml"))

	httpClient := &http.Client{}

	for i := 0; i < 5; i++ {
		res, err := httpClient.Get(m.URL("test"))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode, i)
	}

	res, err := httpClient.Get(m.URL("/test"))
	require.NoError(t, err)
	require.Equal(t, dzhttp.StatusNoMatch, res.StatusCode)
}
