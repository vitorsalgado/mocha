package test

import (
	"crypto/tls"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
	"github.com/vitorsalgado/mocha/to"
)

func TestTLS(t *testing.T) {
	m := mocha.ForTest(t)
	m.StartTLS()

	// allow insecure https request
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	scoped := m.Mock(mocha.Get(to.HaveURLPath("/test")).
		Header("test", to.Equal("hello")).
		Reply(reply.OK()))

	req := testutil.Get(m.Server.URL + "/test")
	req.Header("test", "hello")

	res, err := req.Do()
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	assert.True(t, scoped.IsDone())
}
