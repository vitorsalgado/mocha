package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
)

func TestEnabledDisabledMocks(t *testing.T) {
	httpClient := &http.Client{}
	m := mhttp2.NewAPIWithT(t, mhttp2.Setup().MockFilePatterns("testdata/enabled/*.yaml"))
	m.MustStart()

	testCases := []struct {
		path   string
		status int
	}{
		{"/test", 200},
		{"/hello", 200},
		{"/dev", mhttp2.StatusNoMatch},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			res, err := httpClient.Get(m.URL(tc.path))

			require.NoError(t, err)
			require.Equal(t, tc.status, res.StatusCode)
		})
	}
}
