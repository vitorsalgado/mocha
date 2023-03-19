package test

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestProxy(t *testing.T) {
	proxySrv := NewAPI(Setup().Proxy()).CloseWithT(t)
	proxySrv.MustStart()
	proxyScope := proxySrv.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	targetSrv := NewAPI().CloseWithT(t)
	targetSrv.MustStart()
	targetScope := targetSrv.MustMock(Get(URLPath("/other")).Reply(Created()))

	// client that acts like a browser proxying requests to our server
	proxyURL, _ := url.Parse(proxySrv.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	res, err := client.Get(targetSrv.URL() + "/test")
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(targetSrv.URL() + "/other")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	proxyScope.AssertCalled(t)
	targetScope.AssertCalled(t)
}

func TestProxyTLS(t *testing.T) {
	proxySrv := NewAPI(Setup().Proxy(&ProxyConfig{SSLVerify: false})).CloseWithT(t)
	proxySrv.MustStartTLS()
	proxyScope := proxySrv.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	targetSrv := NewAPI().CloseWithT(t)
	targetSrv.MustStart()
	targetScope := targetSrv.MustMock(Get(URLPath("/other")).Reply(Created()))

	// client that acts like a browser proxying requests to our server
	proxyURL, _ := url.Parse(proxySrv.URL())
	client := &http.Client{Transport: &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	res, err := client.Get(targetSrv.URL() + "/test")
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(targetSrv.URL() + "/other")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	proxyScope.AssertCalled(t)
	targetScope.AssertCalled(t)
}

func TestProxyTLS_CustomCert(t *testing.T) {
	caCert, err := os.ReadFile(filepath.Clean(_tlsCertFile))
	require.NoError(t, err)

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(_tlsCertFile, _tlsKeyFile)
	require.NoError(t, err)

	transport := &http.Transport{TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}}}

	proxySrv := NewAPI(
		Setup().
			Proxy(&ProxyConfig{SSLVerify: true, Transport: transport}).TLSCertKeyPair(_tlsCertFile, _tlsKeyFile)).
		CloseWithT(t)
	proxySrv.MustStartTLS()
	proxyScope := proxySrv.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	targetSrv := NewAPI(Setup().TLSCertKeyPair(_tlsCertFile, _tlsKeyFile)).CloseWithT(t)
	targetSrv.MustStart()
	targetScope := targetSrv.MustMock(Get(URLPath("/other")).Reply(Created()))

	// client that acts like a browser proxying requests to our server
	proxyURL, _ := url.Parse(proxySrv.URL())
	client := &http.Client{Transport: &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}}

	res, err := client.Get(targetSrv.URL() + "/test")
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(targetSrv.URL() + "/other")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	proxyScope.AssertCalled(t)
	targetScope.AssertCalled(t)
}

func TestProxyViaAnotherProxy(t *testing.T) {
	p := NewAPI(Setup().Proxy()).CloseWithT(t)
	p.MustStart()
	scope1 := p.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	v := NewAPI(Setup().Proxy(&ProxyConfig{Via: p.URL()})).CloseWithT(t)
	v.MustStart()
	scope2 := v.MustMock(Get(URLPath("/unknown")).Reply(NoContent()))

	m := NewAPI().CloseWithT(t)
	m.MustStart()
	m.MustMock(Get(URLPath("/other")).Reply(Created()))

	u, _ := url.Parse(v.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}

	res, err := client.Get(m.URL() + "/test")
	require.NoError(t, err)
	scope1.AssertCalled(t)
	scope2.AssertNotCalled(t)
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(m.URL() + "/other")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
}
