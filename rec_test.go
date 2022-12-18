package mocha

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dir := t.TempDir()
	p := New(t, Configure().Name("proxy").Proxy().Record(&RecordConfig{SaveDir: dir, Save: true, SaveBodyToFile: true}))
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

	<-ctx.Done()

	entries, err := os.ReadDir(dir)

	assert.NoError(t, err)
	assert.Len(t, entries, 1)

	p.Close()
	m.Close()
}
