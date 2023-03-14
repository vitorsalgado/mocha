package test

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

const (
	_tlsCertFile = "testdata/tls/srv.crt"
	_tlsKeyFile  = "testdata/tls/srv.key"
)

func TestTLS(t *testing.T) {
	cert, err := tls.LoadX509KeyPair(_tlsCertFile, _tlsKeyFile)
	require.NoError(t, err)

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	testCases := []struct {
		name   string
		config Configurer
	}{
		{"basic", Setup()},
		{"custom key/pair", Setup().TLSCertificateAndKey(_tlsCertFile, _tlsKeyFile)},
		{"tls config", Setup().TLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewT(t, tc.config)
			m.MustStartTLS()

			scoped := m.MustMock(
				Get(URLPath("/test")).
					Header("test", StrictEqual("hello")).
					Reply(OK().PlainText("hello+world")))

			req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
			req.Header.Add("test", "hello")

			res, err := client.Do(req)
			require.NoError(t, err)

			b, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, "hello+world", string(b))
			require.True(t, scoped.HasBeenCalled())
		})
	}
}

func TestTLS_FileSetup(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	m := NewT(t, UseConfig("testdata/tls/tls_app.yaml"))
	m.MustStartTLS()

	scoped := m.MustMock(
		Get(URLPath("/test")).
			Header("test", StrictEqual("hello")).
			Reply(OK().PlainText("hello+world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	req.Header.Add("test", "hello")

	res, err := client.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, "hello+world", string(b))
	require.True(t, scoped.HasBeenCalled())
}

func TestTLSRootCA(t *testing.T) {
	caCert, err := os.ReadFile(filepath.Clean(_tlsCertFile))
	require.NoError(t, err)

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool}},
	}

	target := NewT(t, Setup().TLSCertificateAndKey(_tlsCertFile, _tlsKeyFile))
	target.MustStartTLS()
	target.MustMock(Getf("/test").Reply(Accepted().PlainText("accepted")))

	m := NewT(t, Setup().TLSRootCAs(pool).TLSCertificateAndKey(_tlsCertFile, _tlsKeyFile))
	m.MustStartTLS()
	m.MustMock(Get(URLPath("/test")).
		Header("test", StrictEqual("hello")).
		Reply(From(target.URL()).SSLVerify(true)))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	req.Header.Add("test", "hello")

	res, err := client.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, "accepted", string(b))
	require.Equal(t, http.StatusAccepted, res.StatusCode)
}

func TestTLSRootCA_FileSetup(t *testing.T) {
	caCert, err := os.ReadFile(filepath.Clean(_tlsCertFile))
	require.NoError(t, err)

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool}},
	}

	target := NewT(t, UseConfig("testdata/tls/tls_ca_target.yaml"))
	target.MustStartTLS()
	target.MustMock(Getf("/test").Reply(Accepted().PlainText("accepted")))

	m := NewT(t, UseConfig("testdata/tls/tls_ca.yaml"))
	m.MustStartTLS()
	m.MustMock(Get(URLPath("/test")).
		Header("test", StrictEqual("hello")).
		Reply(From(target.URL()).SSLVerify(true)))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	req.Header.Add("test", "hello")

	res, err := client.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, "accepted", string(b))
	require.Equal(t, http.StatusAccepted, res.StatusCode)
}
