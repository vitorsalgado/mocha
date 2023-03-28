package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
)

func TestSimilarRequestURLPaths(t *testing.T) {
	m := mhttp2.NewAPI()
	m.MustStart()

	defer m.Close()

	scope := m.MustMock(
		mhttp2.Getf("/customers").Reply(mhttp2.OK()),
		mhttp2.Getf("/customers/100").Reply(mhttp2.Accepted()),
		mhttp2.Getf("/customers/100/orders").Reply(mhttp2.Unauthorized()),
		mhttp2.Getf("/customers/100/orders/BR-500").Reply(mhttp2.BadRequest()),
		mhttp2.Getf("/customers/100/orders/BR-500/items").Reply(mhttp2.InternalServerError()))

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
	require.Equal(t, mhttp2.StatusNoMatch, res.StatusCode)

	scope.AssertNumberOfCalls(t, 5)
}
