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
	p := New(t, Configure().Proxy()).CloseOnT(t)
	p.MustStart()
	p.MustMock(Get(URLPath("/test")).Reply(reply.Accepted()))

	m := New(t).CloseOnT(t)
	m.MustStart()
	m.MustMock(Get(URLPath("/other")).Reply(reply.Created()))

	u, _ := url.Parse(p.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}

	res, err := client.Get(m.URL() + "/test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)

	res, err = client.Get(m.URL() + "/other")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}
