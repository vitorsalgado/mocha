package mocha

import (
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestRecording_WithWebProxy(t *testing.T) {
	dir := t.TempDir()
	p := New(t, Configure().Proxy().Record(WithRecordDir(dir)))
	p.MustStart()
	scope1 := p.MustMock(Get(URLPath("/test")).Reply(reply.Accepted()))

	m := New(t)
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

	time.Sleep(1 * time.Second)

	entries, err := os.ReadDir(dir)

	assert.NoError(t, err)
	assert.Len(t, entries, 2)

	p.Close()
	m.Close()
}
