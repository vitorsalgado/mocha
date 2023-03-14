package test

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"

	. "github.com/vitorsalgado/mocha/v3"
)

func TestHTTP2_SettingCAs(t *testing.T) {
	caCert, err := os.ReadFile(_tlsCertFile)
	require.NoError(t, err)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Transport: &http2.Transport{TLSClientConfig: &tls.Config{RootCAs: caCertPool}},
	}

	m := NewT(t, Setup().TLSCertificateAndKey(_tlsCertFile, _tlsKeyFile))
	m.MustStartTLS()
	m.MustMock(Getf("/test").Reply(OK().PlainText("hello world")))

	res, err := client.Get(m.URL("/test"))
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello world", string(b))
	require.Equal(t, "HTTP/2.0", res.Proto)
}
