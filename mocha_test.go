package mocha

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"

	"github.com/vitorsalgado/mocha/v3/internal/testmocks"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(
		m,
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
	)
}

func TestMocha(t *testing.T) {
	m := New(t)
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(
		Get(URLPath("/test")).
			Header("test", Equal("hello")).
			Query("filter", Equal("all")).
			Reply(reply.
				Created().
				BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)
	req.Header.Add("test", "hello")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.True(t, scoped.Called())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}

func TestMocha_NewBasic(t *testing.T) {
	m := NewBasic()
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(
		Get(URLPath("/test")).
			Reply(reply.
				Created().
				BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.True(t, scoped.Called())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}

func TestResponseMapper(t *testing.T) {
	m := New(t)
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(Get(URLPath("/test")).
		Reply(reply.
			OK().
			Map(func(r *reply.Response, rma *reply.MapperArgs) error {
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

	defer m.Close()

	start := time.Now()
	delay := time.Duration(1250) * time.Millisecond

	scoped := m.AddMocks(Get(URLPath("/test")).
		Delay(delay).
		Reply(reply.OK()))

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

func TestErrors(t *testing.T) {
	m := New(t)
	m.Start()

	defer m.Close()

	t.Run("should log errors on reply", func(t *testing.T) {
		scoped := m.AddMocks(Get(URLPath("/test1")).
			ReplyFunc(func(_ http.ResponseWriter, r *http.Request) (*reply.Response, error) {
				return nil, fmt.Errorf("failed to build a response")
			}))

		res, err := testutil.Get(fmt.Sprintf("%s/test1", m.URL())).Do()

		assert.Nil(t, err)
		assert.False(t, scoped.Called())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
	})

	t.Run("should log errors from matchers", func(t *testing.T) {
		scoped := m.AddMocks(Get(URLPath("/test2")).
			Header("test", Func(
				func(_ any) (bool, error) {
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

	defer m.Close()

	fakeT := testmocks.NewFakeNotifier()

	scoped := m.AddMocks(
		Get(URLPath("/test-ok")).
			Reply(reply.OK()))

	assert.Equal(t, 0, scoped.Hits())
	assert.False(t, m.AssertCalled(fakeT))
	assert.True(t, m.AssertNotCalled(fakeT))
	assert.True(t, m.AssertHits(fakeT, 0))
	assert.Equal(t, 0, m.Hits())

	res, err := testutil.Get(m.URL() + "/test-ok").Do()

	assert.NoError(t, err)
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

	defer m.Close()

	m.AddMocks(
		Get(URLPath("/test-1")).
			Reply(reply.OK()),
		Get(URLPath("/test-2")).
			Reply(reply.OK()))

	res, err := testutil.Get(m.URL() + "/test-1").Do()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// disable all mocks
	// should return tea pot for all
	m.Disable()

	res, err = testutil.Get(m.URL() + "/test-1").Do()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, res.StatusCode)

	res, err = testutil.Get(m.URL() + "/test-2").Do()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, res.StatusCode)

	// re-enable mocks again
	m.Enable()

	res, err = testutil.Get(m.URL() + "/test-1").Do()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = testutil.Get(m.URL() + "/test-2").Do()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestMocha_ReplyJust(t *testing.T) {
	t.Run("should return status set on first parameter", func(t *testing.T) {
		m := New(t)
		m.Start()

		defer m.Close()

		scoped := m.AddMocks(
			Post(URLPath("/test")).
				ReplyJust(http.StatusCreated, reply.New().Header("test", "ok")))

		req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run("should overwrite status", func(t *testing.T) {
		m := New(t)
		m.Start()

		defer m.Close()

		scoped := m.AddMocks(
			Post(URLPath("/test")).
				ReplyJust(http.StatusCreated, reply.OK().Header("test", "ok")))

		req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}

type FakeEvents struct{ mock.Mock }

func (h *FakeEvents) OnRequest(e any) {
	h.Called(e)
}

func (h *FakeEvents) OnRequestMatched(e any) {
	h.Called(e)
}

func (h *FakeEvents) OnRequestNotMatched(e any) {
	h.Called(e)
}

func (h *FakeEvents) OnError(e any) {
	h.Called(e)
}

func TestMocha_Subscribe(t *testing.T) {
	f := &FakeEvents{}
	f.On("OnRequest", mock.Anything).Return()
	f.On("OnRequestMatched", mock.Anything).Return()

	m := New(t, Configure().LogLevel(LogSilently).Build()).CloseOnCleanup(t)
	m.Subscribe(EventOnRequest, f.OnRequest)
	m.Subscribe(EventOnRequestMatched, f.OnRequestMatched)
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(
		Get(URLPath("/test")).
			Reply(reply.OK()))

	res, err := testutil.Get(m.URL() + "/test").Do()
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, scoped.Called())

	f.AssertExpectations(t)
}

func TestMocha_Silently(t *testing.T) {
	m := New(t, Configure().LogLevel(LogSilently).Build())
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(
		Get(URLPath("/test")).
			Reply(reply.
				Created().
				BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, scoped.Called())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}

func TestMocha_MatcherCompositions(t *testing.T) {
	m := New(t)
	m.Start()

	defer m.Close()

	scoped := m.AddMocks(
		Get(URLPath("/test")).
			Header("test", Should(Be(Equal("hello")))).
			Query("filter", Is(Equal("all"))).
			Reply(reply.
				Created().
				BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)
	req.Header.Add("test", "hello")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.True(t, scoped.Called())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}
