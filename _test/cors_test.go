package test

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
)

func TestCORS(t *testing.T) {
	m := mocha.New(t, mocha.Configure().CORS().Build())
	m.Start()

	m.Mock(mocha.Get(expect.URLPath("/test")).
		Reply(reply.OK()))

	corsReq := testutil.Request(http.MethodOptions, m.URL()+"/test")
	res, err := corsReq.Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	req := testutil.Request(http.MethodGet, m.URL()+"/test")
	res, err = req.Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
