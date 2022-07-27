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
	"github.com/stretchr/testify/mock"

	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/hooks"
	"github.com/vitorsalgado/mocha/internal/testmocks"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/params"
	"github.com/vitorsalgado/mocha/reply"
)

func TestMocha(t *testing.T) {
	m := New(t)
	m.Start()

	scoped := m.AddMocks(
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
}

func TestMocha_NewBasic(t *testing.T) {
	m := NewBasic()
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(
		Get(expect.URLPath("/test")).
			Reply(reply.
				Created().
				BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
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
}

func TestMocha_Parameters(t *testing.T) {
	key := "k"
	expected := "test"

	m := New(t)
	m.Start()
	m.Parameters().Set(key, expected)

	scoped := m.AddMocks(Get(expect.URLPath("/test")).
		RequestMatches(expect.Func(func(v any, params expect.Args) (bool, error) {
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

	scoped := m.AddMocks(Get(expect.URLPath("/test")).
		Reply(reply.
			OK().
			Map(func(r *reply.Response, rma reply.ResponseMapperArgs) error {
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

func TestResponseDelay(t *testing.T) {
	m := New(t)
	m.Start()

	start := time.Now()
	delay := time.Duration(1250) * time.Millisecond

	scoped := m.AddMocks(Get(expect.URLPath("/test")).
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

	scoped := m.AddMocks(
		Request().
			MatchAfter(expect.Repeat(2)).
			Method(http.MethodGet).
			URL(expect.URLPath("/test")).
			Reply(reply.
				OK()))

	_, _ = testutil.Get(fmt.Sprintf("%s/other", m.URL())).Do()
	_, _ = testutil.Get(fmt.Sprintf("%s/other", m.URL())).Do()

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
	m := New(t)
	m.Start()

	defer m.Close()

	t.Run("should log errors on reply", func(t *testing.T) {
		scoped := m.AddMocks(Get(expect.URLPath("/test1")).
			ReplyFunction(func(r *http.Request, m reply.M, p params.P) (*reply.Response, error) {
				return nil, fmt.Errorf("failed to build a response")
			}))

		res, err := testutil.Get(fmt.Sprintf("%s/test1", m.URL())).Do()

		assert.Nil(t, err)
		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
	})

	t.Run("should log errors from matchers", func(t *testing.T) {
		scoped := m.AddMocks(Get(expect.URLPath("/test2")).
			Header("test", expect.Func(
				func(_ any, _ expect.Args) (bool, error) {
					return false, fmt.Errorf("failed")
				})))

		res, err := testutil.Get(fmt.Sprintf("%s/test2", m.URL())).Do()

		assert.Nil(t, err)
		assert.False(t, scoped.Called())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
	})
}

func TestMocha_Assertions(t *testing.T) {
	m := New(t)
	m.Start()

	fakeT := testmocks.NewFakeNotifier()

	scoped := m.AddMocks(
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

	m.AddMocks(
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

func TestMocha_ReplyJust(t *testing.T) {
	t.Run("should return status set on first parameter", func(t *testing.T) {
		m := New(t)
		m.Start()

		scoped := m.AddMocks(
			Post(expect.URLPath("/test")).
				ReplyJust(http.StatusCreated, reply.New().Header("test", "ok")))

		req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run("should overwrite status", func(t *testing.T) {
		m := New(t)
		m.Start()

		scoped := m.AddMocks(
			Post(expect.URLPath("/test")).
				ReplyJust(http.StatusCreated, reply.OK().Header("test", "ok")))

		req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}

func TestMocha_Context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	m := New(t, Configure().Context(ctx).Build())
	m.Start()

	scoped := m.AddMocks(
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

type FakeEvents struct{ mock.Mock }

func (h *FakeEvents) OnRequest(e hooks.OnRequest) {
	h.Called(e)
}

func (h *FakeEvents) OnRequestMatched(e hooks.OnRequestMatch) {
	h.Called(e)
}

func (h *FakeEvents) OnRequestNotMatched(e hooks.OnRequestNotMatched) {
	h.Called(e)
}

func (h *FakeEvents) OnError(e hooks.OnError) {
	h.Called(e)
}

func TestMocha_Subscribe(t *testing.T) {
	f := &FakeEvents{}
	f.On("OnRequest", mock.AnythingOfType("OnRequest")).Return()
	f.On("OnRequestMatched", mock.Anything).Return()

	m := New(t, Configure().Build())
	m.Subscribe(f)
	m.Start()

	scoped := m.AddMocks(
		Get(expect.URLPath("/test")).
			Reply(reply.OK()))

	res, err := testutil.Get(m.URL() + "/test").Do()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, scoped.Called())

	f.AssertExpectations(t)
}

func TestMocha_Silently(t *testing.T) {
	m := New(t, Configure().LogVerbosity(LogSilently).Build())
	m.Start()

	scoped := m.AddMocks(
		Get(expect.URLPath("/test")).
			Reply(reply.
				Created().
				BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)

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
}
