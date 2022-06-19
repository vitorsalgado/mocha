package test

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/internal/assert"
)

func TestFormUrlEncoded(t *testing.T) {
	m := mocha.NewT(t)

	scoped := m.Mock(mocha.Post(mocha.URLPath("/test")).
		Reply(mocha.OK()))

	data := url.Values{}
	data.Set("var1", "dev")
	data.Set("vqr2", "qa")

	req, _ := http.NewRequest(http.MethodPost, m.Server.URL+"/test", strings.NewReader(data.Encode()))
	req.Header.Add("test", "hello")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	assert.True(t, scoped.IsDone())
}
