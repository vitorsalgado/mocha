package test

import (
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

func TestExternalResponseBodies(t *testing.T) {
	m := mocha.NewAPI()
	m.MustStart()

	defer m.Close()

	client := &http.Client{}
	cases := []struct {
		mock *dzhttp.HTTPMockBuilder
		body string
	}{
		{m.Getf("/test").
			Reply(dzhttp.OK().BodyFile("./testdata/external_body/01.txt")), "hello world\n"},
		{m.Getf("/test").
			Reply(dzhttp.OK().Gzip().BodyFile("./testdata/external_body/01.txt")), "hello world\n"},
	}

	for i, tc := range cases {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m.MustMock(tc.mock)

			req, _ := http.NewRequest(http.MethodGet, m.URL("/test"), nil)
			res, err := client.Do(req)
			require.NoError(t, err)

			b, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, tc.body, string(b))

			m.Clean()
		})
	}
}
