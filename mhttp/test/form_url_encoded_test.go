package test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/mhttp"
	"github.com/vitorsalgado/mocha/v3/mhttpv"
)

func TestFormUrlEncoded(t *testing.T) {
	httpClient := &http.Client{}
	m := mhttp.NewAPIWithT(t)
	m.MustStart()

	scoped := m.MustMock(mhttp.Post(matcher.URLPath("/test")).
		FormField("var1", matcher.StrictEqual("dev")).
		FormField("var2", matcher.Contain("q")).
		Reply(mhttp.OK()))

	data := url.Values{}
	data.Set("var1", "dev")
	data.Set("var2", "qa")

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader(data.Encode()))
	req.Header.Add("test", "hello")
	req.Header.Add(mhttpv.HeaderContentType, mhttpv.MIMEFormURLEncoded)
	res, err := httpClient.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.True(t, scoped.HasBeenCalled())
}

func TestFormUrlEncoded_InvalidFieldValues(t *testing.T) {
	m := mhttp.NewAPIWithT(t)
	scoped, err := m.Mock(mhttp.FromFile("testdata/form_url_encoded/01_invalid.yaml"))

	require.Nil(t, scoped)
	require.Error(t, err)
}

func TestFormUrlEncoded_FromFileMock(t *testing.T) {
	httpClient := &http.Client{}
	m := mhttp.NewAPIWithT(t)
	m.MustStart()

	scoped := m.MustMock(mhttp.FromFile("testdata/form_url_encoded/02_valid.yaml"))

	testCases := []struct {
		name           string
		data           func() url.Values
		expectedStatus int
	}{
		{"match", func() url.Values {
			data := url.Values{}
			data.Set("name", "nice name")
			data.Set("address", "berlin+germany")
			data.Set("active", "true")
			data.Set("live", "false")
			data.Set("money", "2550.50")
			data.Set("code", "10")
			data.Set("job", "dev")

			return data
		}, http.StatusNoContent},

		{"no match (not equals)", func() url.Values {
			data := url.Values{}
			data.Set("name", "nice name")
			data.Set("address", "berlin+germany")
			data.Set("active", "true")
			data.Set("live", "false")
			data.Set("money", "2550.50")
			data.Set("code", "10")
			data.Set("job", "qa")

			return data
		}, mhttp.StatusNoMatch},

		{"no match (missing +)", func() url.Values {
			data := url.Values{}
			data.Set("name", "nice name")
			data.Set("address", "berlin germany")
			data.Set("active", "true")
			data.Set("live", "false")
			data.Set("money", "2550.50")
			data.Set("code", "10")
			data.Set("job", "dev")

			return data
		}, mhttp.StatusNoMatch},

		{"no match (missing field)", func() url.Values {
			data := url.Values{}
			data.Set("name", "nice name")
			data.Set("address", "berlin+germany")
			data.Set("active", "true")
			data.Set("live", "false")
			data.Set("code", "10")
			data.Set("job", "dev")

			return data
		}, mhttp.StatusNoMatch},

		{"match (float)", func() url.Values {
			data := url.Values{}
			data.Set("name", "nice name")
			data.Set("address", "berlin+germany")
			data.Set("active", "true")
			data.Set("live", "false")
			data.Set("money", "2550.50")
			data.Set("code", "10")
			data.Set("job", "dev")

			return data
		}, http.StatusNoContent},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader(tc.data().Encode()))
			req.Header.Add(mhttpv.HeaderContentType, mhttpv.MIMEFormURLEncoded)
			res, err := httpClient.Do(req)

			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}

	scoped.AssertNumberOfCalls(t, 2)
}
