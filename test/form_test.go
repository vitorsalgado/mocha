package test

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/internal/header"
	"github.com/vitorsalgado/mocha/internal/mime"
	"github.com/vitorsalgado/mocha/matchers"
	"github.com/vitorsalgado/mocha/reply"
)

func TestFormUrlEncoded(t *testing.T) {
	m := mocha.ForTest(t)
	m.Start()

	scoped := m.Mock(mocha.Post(matchers.URLPath("/test")).
		Reply(reply.OK()))

	data := url.Values{}
	data.Set("var1", "dev")
	data.Set("vqr2", "qa")

	req, _ := http.NewRequest(http.MethodPost, m.Server.URL+"/test", strings.NewReader(data.Encode()))
	req.Header.Add("test", "hello")
	req.Header.Add(header.ContentType, mime.ContentTypeFormURLEncoded)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	assert.True(t, scoped.IsDone())
}
