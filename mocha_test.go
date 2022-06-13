package mocha

import (
	"github.com/vitorsalgado/mocha/internal/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestMocha(t *testing.T) {
	t.Run("should mock request", func(t *testing.T) {
		m := NewT(t)

		scoped := m.Mock(Get(URLPath("/test")).
			Header("test", Equal("hello")).
			Query("filter", Equal("all")).
			Reply(Created().BodyStr("hello world")))

		req, _ := http.NewRequest(http.MethodGet, m.Server.URL+"/test?filter=all", nil)
		req.Header.Add("test", "hello")
		req.URL.Query().Add("filter", "all")

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
