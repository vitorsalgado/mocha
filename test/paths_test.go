package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
)

func TestSimilarRequestURLPaths(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	scope := m.MustMock(
		mocha.Getf("/customers").Reply(mocha.OK()),
		mocha.Getf("/customers/100").Reply(mocha.Accepted()),
		mocha.Getf("/customers/100/orders").Reply(mocha.Unauthorized()),
		mocha.Getf("/customers/100/orders/BR-500").Reply(mocha.BadRequest()),
		mocha.Getf("/customers/100/orders/BR-500/items").Reply(mocha.InternalServerError()))

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
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)

	scope.AssertNumberOfCalls(t, 5)
}
