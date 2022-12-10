package mocha

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestProxy(t *testing.T) {
	p := New(t, Configure().Proxy()).CloseWithT(t)
	p.MustStart()
	scope1 := p.MustMock(Get(URLPath("/test")).Reply(reply.Accepted()))

	m := New(t).CloseWithT(t)
	m.MustStart()
	scope2 := m.MustMock(Get(URLPath("/other")).Reply(reply.Created()))

	u, _ := url.Parse(p.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}

	res, err := client.Get(m.URL() + "/test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(m.URL() + "/other")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	scope1.AssertCalled(t)
	scope2.AssertCalled(t)
}

func TestProxy_ViaProxy(t *testing.T) {
	p := New(t, WithProxy()).CloseWithT(t)
	p.MustStart()
	scope1 := p.MustMock(Get(URLPath("/test")).Reply(reply.Accepted()))

	u1, _ := url.Parse(p.URL())
	v := New(t, WithProxy(&ProxyConfig{ProxyVia: u1})).CloseWithT(t)
	v.MustStart()
	scope2 := v.MustMock(Get(URLPath("/unknown")).Reply(reply.NoContent()))

	m := New(t).CloseWithT(t)
	m.MustStart()
	m.MustMock(Get(URLPath("/other")).Reply(reply.Created()))

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
