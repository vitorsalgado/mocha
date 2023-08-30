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

	. "github.com/vitorsalgado/mocha/v3/dzhttp"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

const (
	_tlsCertFile       = "testdata/cert/cert.pem"
	_tlsKeyFile        = "testdata/cert/key.pem"
	_tlsClientCertFile = "testdata/cert/cert_client.pem"
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
		{"default", Setup()},
		{"custom key/pair", Setup().TLSCertKeyPair(_tlsCertFile, _tlsKeyFile)},
		{"custom tls config", Setup().TLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}})},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m := NewAPI(tc.config).CloseWithT(t)
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
	m := NewAPI(UseConfig("testdata/tls/tls_app.yaml")).CloseWithT(t)
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

func TestTLSMutual(t *testing.T) {
	cert, err := tls.LoadX509KeyPair(_tlsCertFile, _tlsKeyFile)
	require.NoError(t, err)

	caCert, err := os.ReadFile(filepath.Clean(_tlsClientCertFile))
	require.NoError(t, err)

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
				ClientAuth:   tls.RequireAndVerifyClientCert,
			}},
	}

	target := NewAPI(Setup().TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile)).CloseWithT(t)
	target.MustStartTLS()
	target.MustMock(Getf("/test").Reply(Accepted().PlainText("accepted")))

	req, _ := http.NewRequest(http.MethodGet, target.URL("/test"), nil)
	res, err := client.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusAccepted, res.StatusCode)
	require.Equal(t, "accepted", string(b))
}

func TestTLSMutualWithProxy(t *testing.T) {
	cert, err := tls.LoadX509KeyPair(_tlsCertFile, _tlsKeyFile)
	require.NoError(t, err)

	caCert, err := os.ReadFile(filepath.Clean(_tlsClientCertFile))
	require.NoError(t, err)

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
			}},
	}

	target := NewAPI(Setup().TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile)).CloseWithT(t)
	target.MustStartTLS()
	target.MustMock(Getf("/test").Reply(Accepted().PlainText("accepted")))

	m := NewAPI(Setup().TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile)).CloseWithT(t)
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
	require.Equal(t, http.StatusAccepted, res.StatusCode)
	require.Equal(t, "accepted", string(b))
}

func TestTLSMutualWithProxy_FileSetup(t *testing.T) {
	cert, err := tls.LoadX509KeyPair(_tlsCertFile, _tlsKeyFile)
	require.NoError(t, err)

	caCert, err := os.ReadFile(filepath.Clean(_tlsClientCertFile))
	require.NoError(t, err)

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
			}},
	}

	target := NewAPI(UseConfig("testdata/tls/tls_ca_target.yaml")).CloseWithT(t)
	target.MustStartTLS()
	target.MustMock(Getf("/test").Reply(Accepted().PlainText("accepted")))

	m := NewAPI(UseConfig("testdata/tls/tls_ca.yaml")).CloseWithT(t)
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
