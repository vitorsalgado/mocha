package dzhttp

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
)

func TestReplyEcho(t *testing.T) {
	m := NewAPI().CloseWithT(t)
	m.MustStart()

	defer m.Close()

	expectedBody := `{"status":"ok"}`
	scope := m.MustMock(Postf("/test").Reply(Echo().Log()))
	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, m.URL("/test"), strings.NewReader(expectedBody))
	req.Header.Add(httpval.HeaderContentType, httpval.MIMETextPlainCharsetUTF8)

	res, err := httpClient.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.NoError(t, res.Body.Close())
	require.True(t, scope.AssertCalled(t))
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, httpval.MIMETextPlainCharsetUTF8, res.Header.Get(httpval.HeaderContentType))
	require.Contains(t, string(b), expectedBody)
	require.Contains(t, string(b), res.Header.Get("x-custom-header"))
}
