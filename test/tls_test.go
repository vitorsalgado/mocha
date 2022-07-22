package test

import (
	"crypto/tls"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
)

func TestTLS(t *testing.T) {
	m := mocha.New(t)
	m.StartTLS()

	// allow insecure https request
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	scoped := m.AddMocks(mocha.Get(expect.URLPath("/test")).
		Header("test", expect.ToEqual("hello")).
		Reply(reply.OK()))

	req := testutil.Get(m.URL() + "/test")
	req.Header("test", "hello")

	res, err := req.Do()
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	assert.True(t, scoped.Called())
}
