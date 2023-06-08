package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/httpd"
)

func TestSimilarRequestURLPaths(t *testing.T) {
	m := httpd.NewAPI()
	m.MustStart()

	defer m.Close()

	scope := m.MustMock(
		httpd.Getf("/customers").Reply(httpd.OK()),
		httpd.Getf("/customers/100").Reply(httpd.Accepted()),
		httpd.Getf("/customers/100/orders").Reply(httpd.Unauthorized()),
		httpd.Getf("/customers/100/orders/BR-500").Reply(httpd.BadRequest()),
		httpd.Getf("/customers/100/orders/BR-500/items").Reply(httpd.InternalServerError()))

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
	require.Equal(t, httpd.StatusNoMatch, res.StatusCode)

	scope.AssertNumberOfCalls(t, 5)
}
