package mocha

import (
	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/reply"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/matcher"
)

func TestMocha(t *testing.T) {
	t.Run("should mock request", func(t *testing.T) {
		m := ForTest(t)
		m.Start()

		scoped := m.Mock(Get(matcher.URLPath("/test")).
			Header("test", matcher.EqualTo("hello")).
			Query("filter", matcher.EqualTo("all")).
			Reply(reply.Created().BodyString("hello world")))

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

type J struct {
	Name string `json:"name"`
	OK   bool   `json:"ok"`
}

func TestPostJSON(t *testing.T) {
	m := ForTest(t)
	m.Start()

	scoped := m.Mock(Post(matcher.URLPath("/test")).
		Header("test", matcher.EqualTo("hello")).
		Body(
			matcher.JSONPath("name", matcher.Equal("dev")), matcher.JSONPath("ok", matcher.Equal(true))).
		Reply(reply.OK()))

	req := testutil.PostJSON(m.Server.URL+"/test", &J{Name: "dev", OK: true})
	req.Header("test", "hello")

	res, err := req.Do()
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	assert.True(t, scoped.IsDone())
}
