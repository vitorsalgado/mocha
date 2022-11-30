package test

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestTLS(t *testing.T) {
	m := mocha.New(t)
	m.StartTLS()

	defer m.Close()

	// allow insecure https request
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	scoped := m.AddMocks(mocha.Get(expect.URLPath("/test")).
		Header("test", expect.ToEqual("hello")).
		Reply(reply.OK()))

	req := testutil.Get(m.URL() + "/test")
	req.Header("test", "hello")

	res, err := req.Do()
	assert.NoError(t, err)

	defer res.Body.Close()

	assert.True(t, scoped.Called())
}
