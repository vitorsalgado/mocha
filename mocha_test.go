package mocha

import (
	"context"
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

func TestHTTPMethods(t *testing.T) {
	m := New(t)
	m.Start()

	t.Run("should mock GET", func(t *testing.T) {
		scoped := m.Mock(
			Get(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 1, scoped.Hits())

		other, err := testutil.Post(m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
	})

	t.Run("should mock POST", func(t *testing.T) {
		scoped := m.Mock(
			Post(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.Post(m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
	})

	t.Run("should mock PUT", func(t *testing.T) {
		scoped := m.Mock(
			Put(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.NewRequest(http.MethodPut, m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
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

func TestMocha_Parameters(t *testing.T) {
	key := "k"
	expected := "test"

	m := New(t)
	m.Start()
	m.Parameters().Set(key, expected)

	scoped := m.Mock(Get(expect.URLPath("/test")).
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

	scoped.AssertCalled(t)
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

	scoped.AssertCalled(t)
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

	scoped.AssertCalled(t)
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

	scoped.AssertCalled(t)
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
		fake.AssertNumberOfCalls(t, "Logf", 1)
	})

	t.Run("should log errors from matchers", func(t *testing.T) {
		scoped := m.Mock(Get(expect.URLPath("/test2")).
			Header("test", expect.Func(
				func(_ any, _ expect.Args) (bool, error) {
					return false, fmt.Errorf("failed")
				})))

		res, err := testutil.Get(fmt.Sprintf("%s/test2", m.URL())).Do()

		assert.Nil(t, err)
		assert.False(t, scoped.Called())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		fake.AssertNumberOfCalls(t, "Logf", 2)
	})
}

func TestMocha_Assertions(t *testing.T) {
	m := New(t)
	m.Start()

	fakeT := mocks.NewT()

	scoped := m.Mock(
		Get(expect.URLPath("/test-ok")).
			Reply(reply.OK()))

	assert.Equal(t, 0, scoped.Hits())
	assert.False(t, m.AssertCalled(fakeT))
	assert.True(t, m.AssertNotCalled(fakeT))
	assert.True(t, m.AssertHits(fakeT, 0))
	assert.Equal(t, 0, m.Hits())

	res, err := testutil.Get(m.URL() + "/test-ok").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, scoped.Hits())
	assert.True(t, m.AssertCalled(fakeT))
	assert.False(t, m.AssertNotCalled(fakeT))
	assert.True(t, m.AssertHits(fakeT, 1))
	assert.Equal(t, 1, m.Hits())
}

func TestMocha_Enable_Disable(t *testing.T) {
	m := New(t)
	m.Start()

	m.Mock(
		Get(expect.URLPath("/test-1")).
			Reply(reply.OK()),
		Get(expect.URLPath("/test-2")).
			Reply(reply.OK()))

	res, err := testutil.Get(m.URL() + "/test-1").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// disable all mocks
	// should return tea pot for all
	m.Disable()

	res, err = testutil.Get(m.URL() + "/test-1").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusTeapot, res.StatusCode)

	res, err = testutil.Get(m.URL() + "/test-2").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusTeapot, res.StatusCode)

	// re-enable mocks again
	m.Enable()

	res, err = testutil.Get(m.URL() + "/test-1").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = testutil.Get(m.URL() + "/test-2").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestMocha_Context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	m := New(Noop(), Configure().Context(ctx).Build())
	m.Start()

	scoped := m.Mock(
		Get(expect.URLPath("/test")).
			Reply(reply.OK()))

	res, err := testutil.Get(m.URL() + "/test").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, scoped.Called())

	cancel()

	res, err = testutil.Get(m.URL() + "/test").Do()

	assert.NotNil(t, err, "server was supposed to be closed")
	assert.Nil(t, res)
}
