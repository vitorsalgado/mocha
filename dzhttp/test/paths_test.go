package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

func TestSimilarRequestURLPaths(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	scope := m.MustMock(
		dzhttp.Getf("/customers").Reply(dzhttp.OK()),
		dzhttp.Getf("/customers/100").Reply(dzhttp.Accepted()),
		dzhttp.Getf("/customers/100/orders").Reply(dzhttp.Unauthorized()),
		dzhttp.Getf("/customers/100/orders/BR-500").Reply(dzhttp.BadRequest()),
		dzhttp.Getf("/customers/100/orders/BR-500/items").Reply(dzhttp.InternalServerError()))

	res, err := http.DefaultClient.Get(m.URL() + "/customers")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.DefaultClient.Get(m.URL() + "/customers/100")
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = http.DefaultClient.Get(m.URL() + "/customers/100/orders")
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	res, err = http.DefaultClient.Get(m.URL() + "/customers/100/orders/BR-500")
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	res, err = http.DefaultClient.Get(m.URL() + "/customers/100/orders/BR-500/items")
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	res, err = http.DefaultClient.Get(m.URL() + "/customers/orders")
	require.NoError(t, err)
	require.Equal(t, dzhttp.StatusNoMatch, res.StatusCode)

	scope.AssertNumberOfCalls(t, 5)
}
