package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/httpd"
)

func TestRepeat(t *testing.T) {
	m := httpd.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(httpd.Get(matcher.URLPath("/test")).
		Times(3).
		Reply(httpd.OK()))

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
	require.Equal(t, httpd.StatusNoMatch, res.StatusCode)
	require.NoError(t, err)
}

func TestRepeat_Once(t *testing.T) {
	m := httpd.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(httpd.Get(matcher.URLPath("/test")).
		Once().
		Reply(httpd.OK()))

	res, _ := http.Get(m.URL() + "/test")
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _ = http.Get(m.URL() + "/test")
	require.Equal(t, httpd.StatusNoMatch, res.StatusCode)
}

func TestRepeat_FileSetup(t *testing.T) {
	m := httpd.NewAPIWithT(t)
	m.MustStart()
	m.MustMock(httpd.FromFile("testdata/repeat/1_repeat.yaml"))

	httpClient := &http.Client{}

	for i := 0; i < 5; i++ {
		res, err := httpClient.Get(m.URL("test"))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode, i)
	}

	res, err := httpClient.Get(m.URL("/test"))
	require.NoError(t, err)
	require.Equal(t, httpd.StatusNoMatch, res.StatusCode)
}
