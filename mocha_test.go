package mocha

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/reply"
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
			matcher.JSONPath("name", matcher.EqualAny("dev")), matcher.JSONPath("ok", matcher.EqualAny(true))).
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

func TestCustomParameters(t *testing.T) {
	key := "k"
	expected := "test"

	m := ForTest(t)
	m.Start()
	m.Parameters().Set(key, expected)

	scope := m.Mock(Get(matcher.URLPath("/test")).
		Matches(func(v any, params matcher.Args) (bool, error) {
			p, _ := params.Params.Get(key)
			return p.(string) == expected, nil
		}).
		Reply(reply.Accepted()))

	req := testutil.Get(fmt.Sprintf("%s/test", m.Server.URL))
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, scope.Done())
	assert.Equal(t, http.StatusAccepted, res.StatusCode)
}

func TestResponseMapper(t *testing.T) {
	m := ForTest(t)
	m.Start()

	scoped := m.Mock(Get(matcher.URLPath("/test")).
		Reply(reply.
			OK().
			Map(func(r *mock.Response, rma mock.ResponseMapperArgs) error {
				r.Header.Add("x-test", rma.Request.Header.Get("x-param"))
				return nil
			})))

	req := testutil.Get(fmt.Sprintf("%s/test", m.Server.URL))
	req.Header("x-param", "dev")

	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, scoped.Done())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "dev", res.Header.Get("x-test"))
}

func TestDelay(t *testing.T) {
	m := ForTest(t)
	m.Start()

	start := time.Now()
	delay := time.Duration(1250) * time.Millisecond

	scoped := m.Mock(Get(matcher.URLPath("/test")).
		Reply(reply.
			OK().
			Delay(delay)))

	req := testutil.Get(fmt.Sprintf("%s/test", m.Server.URL))
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	elapsed := time.Since(start)

	assert.Nil(t, scoped.Done())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.GreaterOrEqual(t, elapsed, delay)
}
