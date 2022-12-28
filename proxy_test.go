package mocha

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestProxy(t *testing.T) {
	proxySrv := New(Configure().Proxy()).CloseWithT(t)
	proxySrv.MustStart()
	proxyScope := proxySrv.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	targetSrv := New().CloseWithT(t)
	targetSrv.MustStart()
	targetScope := targetSrv.MustMock(Get(URLPath("/other")).Reply(Created()))

	// client that acts like a browser proxying requests to our server
	proxyURL, _ := url.Parse(proxySrv.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	res, err := client.Get(targetSrv.URL() + "/test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(targetSrv.URL() + "/other")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	proxyScope.AssertCalled(t)
	targetScope.AssertCalled(t)
}

func TestProxy_ViaProxy(t *testing.T) {
	p := New(WithProxy()).CloseWithT(t)
	p.MustStart()
	scope1 := p.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	v := New(WithProxy(&ProxyConfig{ProxyVia: p.URL()})).CloseWithT(t)
	v.MustStart()
	scope2 := v.MustMock(Get(URLPath("/unknown")).Reply(NoContent()))

	m := New().CloseWithT(t)
	m.MustStart()
	m.MustMock(Get(URLPath("/other")).Reply(Created()))

	u, _ := url.Parse(v.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}

	res, err := client.Get(m.URL() + "/test")
	assert.NoError(t, err)
	scope1.AssertCalled(t)
	scope2.AssertNotCalled(t)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(m.URL() + "/other")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}
