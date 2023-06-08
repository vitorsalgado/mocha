package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/httpd"
)

func TestEnabledDisabledMocks(t *testing.T) {
	httpClient := &http.Client{}
	m := httpd.NewAPIWithT(t, httpd.Setup().MockFilePatterns("testdata/enabled/*.yaml"))
	m.MustStart()

	testCases := []struct {
		path   string
		status int
	}{
		{"/test", 200},
		{"/hello", 200},
		{"/dev", httpd.StatusNoMatch},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			res, err := httpClient.Get(m.URL(tc.path))

			require.NoError(t, err)
			require.Equal(t, tc.status, res.StatusCode)
		})
	}
}
