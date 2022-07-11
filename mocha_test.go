package mocha

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/core/mocks"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/parameters"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
)

type TestModel struct {
	Name string `json:"name"`
	OK   bool   `json:"ok"`
}

func TestMocha(t *testing.T) {
	t.Run("should mock request", func(t *testing.T) {
		m := New(t)
		m.Start()

		scoped := m.Mock(
			Get(expect.URLPath("/test")).
				Header("test", expect.ToEqual("hello")).
				Query("filter", expect.ToEqual("all")).
				Reply(reply.
					Created().
					BodyString("hello world")))

		req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)
		req.Header.Add("test", "hello")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, scoped.Called())
		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, string(body), "hello world")
	})
}

func TestPostJSON(t *testing.T) {
	m := New(t)
	m.Start()

	scoped := m.Mock(Post(expect.URLPath("/test")).
		Header("test", expect.ToEqual("hello")).
		Body(
			expect.JSONPath("name", expect.ToEqual("dev")), expect.JSONPath("ok", expect.ToEqual(true))).
		Reply(reply.OK()))

	req := testutil.PostJSON(m.URL()+"/test", &TestModel{Name: "dev", OK: true})
	req.Header("test", "hello")

	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.True(t, scoped.Called())
}

func TestCustomParameters(t *testing.T) {
	key := "k"
	expected := "test"

	m := New(t)
	m.Start()
	m.Parameters().Set(key, expected)

	scope := m.Mock(Get(expect.URLPath("/test")).
		Matches(expect.Func(func(v any, params expect.Args) (bool, error) {
			p, _ := params.Params.Get(key)
			return p.(string) == expected, nil
		})).
		Reply(reply.Accepted()))

	req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	scope.MustHaveBeenCalled(t)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)
}

func TestResponseMapper(t *testing.T) {
	m := New(t)
	m.Start()

	scoped := m.Mock(Get(expect.URLPath("/test")).
		Reply(reply.
			OK().
			Map(func(r *core.Response, rma core.ResponseMapperArgs) error {
				r.Header.Add("x-test", rma.Request.Header.Get("x-param"))
				return nil
			})))

	req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
	req.Header("x-param", "dev")

	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	scoped.MustHaveBeenCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "dev", res.Header.Get("x-test"))
}

func TestDelay(t *testing.T) {
	m := New(t)
	m.Start()

	start := time.Now()
	delay := time.Duration(1250) * time.Millisecond

	scoped := m.Mock(Get(expect.URLPath("/test")).
		Reply(reply.
			OK().
			Delay(delay)))

	req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	elapsed := time.Since(start)

	scoped.MustHaveBeenCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.GreaterOrEqual(t, elapsed, delay)
}

func TestPostExpectations(t *testing.T) {
	m := New(t)
	m.Start()

	scoped := m.Mock(
		Request().
			MatchAfter(expect.Repeat(2)).
			Method("GET").
			URL(expect.URLPath("/test")).
			Reply(reply.
				OK()))

	testutil.Get(fmt.Sprintf("%s/other", m.URL())).Do()
	testutil.Get(fmt.Sprintf("%s/other", m.URL())).Do()

	res, _ := testutil.Get(fmt.Sprintf("%s/other", m.URL())).Do()
	assert.Equal(t, res.StatusCode, http.StatusTeapot)

	res, _ = testutil.Get(fmt.Sprintf("%s/test", m.URL())).Do()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, _ = testutil.Get(fmt.Sprintf("%s/test", m.URL())).Do()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, _ = testutil.Get(fmt.Sprintf("%s/test", m.URL())).Do()
	assert.Equal(t, res.StatusCode, http.StatusTeapot)

	scoped.MustHaveBeenCalled(t)
}

func TestErrors(t *testing.T) {
	fake := mocks.NewT()

	m := New(fake)
	m.Start()

	defer m.Close()

	t.Run("should log errors on reply", func(t *testing.T) {
		scoped := m.Mock(Get(expect.URLPath("/test1")).
			ReplyFunction(func(r *http.Request, m *core.Mock, p parameters.Params) (*core.Response, error) {
				return nil, fmt.Errorf("failed to build a response")
			}))

		res, err := testutil.Get(fmt.Sprintf("%s/test1", m.URL())).Do()

		assert.Nil(t, err)
		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		fake.AssertNumberOfCalls(t, "Errorf", 1)
	})

	t.Run("should log errors from matchers", func(t *testing.T) {
		scoped := m.Mock(Get(expect.URLPath("/test2")).
			Header("test", expect.Func(
				func(_ string, _ expect.Args) (bool, error) {
					return false, fmt.Errorf("failed")
				})))

		res, err := testutil.Get(fmt.Sprintf("%s/test2", m.URL())).Do()

		assert.Nil(t, err)
		assert.False(t, scoped.Called())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		fake.AssertNumberOfCalls(t, "Errorf", 2)
	})
}

func TestExpect(t *testing.T) {
	m := New(t)
	m.Start()

	scoped := m.Mock(Get(expect.URLPath("/test")).
		Cond(Expect(Header("hello")).ToEqual("world")).
		Reply(reply.
			OK()))

	req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
	req.Header("hello", "world")
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	scoped.MustHaveBeenCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
