package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

func TestRandom_SetupFromFileWithSeed(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(dzhttp.FromFile("testdata/response_random/rand_01.yaml"))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}
