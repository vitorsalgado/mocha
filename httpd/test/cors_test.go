package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/httpd"
)

func TestCORS(t *testing.T) {
	client := &http.Client{}
	m := httpd.NewAPI(httpd.Setup().CORS())
	m.MustStart()

	defer m.Close()

	m.MustMock(httpd.Get(matcher.URLPath("/test")).
		Reply(httpd.OK()))

	corsReq, _ := http.NewRequest(http.MethodOptions, m.URL()+"/test", nil)
	res, err := client.Do(corsReq)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	res, err = client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
