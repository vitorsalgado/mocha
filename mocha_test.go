package mocha

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/vitorsalgado/mocha/internal/assert"
)

func TestMocha(t *testing.T) {
	t.Run("should mock request", func(t *testing.T) {
		m := NewT(t)

		scoped := m.Mock(Get(URLPath("/test")).
			Header("test", EqualTo("hello")).
			Query("filter", EqualTo("all")).
			Reply(Created().BodyStr("hello world")))

		req, _ := http.NewRequest(http.MethodGet, m.Server.URL+"/test?filter=all", nil)
		req.Header.Add("test", "hello")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)

		assert.Nil(t, err)
		assert.True(t, scoped.IsDone())
		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, string(body), "hello world")
	})
}

func TestFormUrlEncoded(t *testing.T) {
	m := NewT(t)

	scoped := m.Mock(Post(URLPath("/test")).
		Reply(OK()))

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

type J struct {
	Name string `json:"name"`
	OK   bool   `json:"ok"`
}

func TestPostJSON(t *testing.T) {
	m := NewT(t)

	scoped := m.Mock(Post(URLPath("/test")).
		Header("test", EqualTo("hello")).
		Body(
			JSONPath("name", Equal("dev")), JSONPath("ok", Equal(true))).
		Reply(OK()))

	body, err := json.Marshal(&J{Name: "dev", OK: true})
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, m.Server.URL+"/test", bytes.NewReader(body))
	req.Header.Add("test", "hello")
	req.Header.Add("content-type", ContentTypeJSON)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	assert.True(t, scoped.IsDone())
}
