package mocha

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"

	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
	"github.com/vitorsalgado/mocha/v3/x/event"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(
		m,
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
	)
}

func TestMocha(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test")).
			Header("test", Equal("hello")).
			Query("filter", Equal("all")).
			Reply(reply.
				Created().
				PlainText("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)
	req.Header.Add("test", "hello")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.True(t, scoped.HasBeenCalled())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}

func TestResponseMapper(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(Get(URLPath("/test")).
		Reply(reply.
			OK()).
		Map(func(r *reply.Stub, rma *MapperIn) error {
			r.Header.Add("x-test", rma.Request.Header.Get("x-param"))
			return nil
		}))

	req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
	req.Header("x-param", "dev")

	res, err := req.Do()
	if err != nil {
		t.Error(err)
	}

	scoped.AssertCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "dev", res.Header.Get("x-test"))
}

func TestResponseDelay(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	start := time.Now()
	delay := 250 * time.Millisecond

	scoped := m.MustMock(Get(URLPath("/test")).
		Delay(delay).
		Reply(reply.OK()))

	req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
	res, err := req.Do()
	if err != nil {
		log.Panic(err)
	}

	elapsed := time.Since(start)

	scoped.AssertCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.GreaterOrEqual(t, elapsed, delay)
}

func TestErrors(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test2")).
			Header("test", Func(
				func(_ any) (bool, error) {
					return false, fmt.Errorf("failed")
				})).
			Reply(reply.OK()))

	res, err := testutil.Get(fmt.Sprintf("%s/test2", m.URL())).Do()

	assert.NoError(t, err)
	assert.False(t, scoped.HasBeenCalled())
	assert.Equal(t, StatusRequestDidNotMatch, res.StatusCode)
}

func TestMocha_Assertions(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	fakeT := NewFakeNotifier()

	scoped := m.MustMock(
		Get(URLPath("/test-ok")).
			Reply(reply.OK()))

	assert.Equal(t, 0, scoped.Hits())
	assert.False(t, m.AssertCalled(fakeT))
	assert.True(t, m.AssertNotCalled(fakeT))
	assert.True(t, m.AssertNumberOfCalls(fakeT, 0))
	assert.Equal(t, 0, m.Hits())

	res, err := testutil.Get(m.URL() + "/test-ok").Do()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, scoped.Hits())
	assert.True(t, m.AssertCalled(fakeT))
	assert.False(t, m.AssertNotCalled(fakeT))
	assert.True(t, m.AssertNumberOfCalls(fakeT, 1))
	assert.Equal(t, 1, m.Hits())
}

func TestMocha_Enable_Disable(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	m.MustMock(
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
	assert.Equal(t, StatusRequestDidNotMatch, res.StatusCode)

	res, err = testutil.Get(m.URL() + "/test-2").Do()

	assert.NoError(t, err)
	assert.Equal(t, StatusRequestDidNotMatch, res.StatusCode)

	// re-enable mocks again
	m.Enable()

	res, err = testutil.Get(m.URL() + "/test-1").Do()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = testutil.Get(m.URL() + "/test-2").Do()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
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

	m := New(Configure().LogLevel(LogSilently)).CloseWithT(t)
	m.MustSubscribe(event.EventOnRequest, f.OnRequest)
	m.MustSubscribe(event.EventOnRequestMatched, f.OnRequestMatched)
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test")).
			Reply(reply.OK()))

	res, err := testutil.Get(m.URL() + "/test").Do()
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, scoped.HasBeenCalled())

	f.AssertExpectations(t)
}

func TestMocha_Silently(t *testing.T) {
	m := New(Configure().LogLevel(LogSilently))
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test")).
			Reply(reply.
				Created().
				PlainText("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, scoped.HasBeenCalled())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}

func TestMocha_MatcherCompositions(t *testing.T) {
	// m := New()
	// m.MustStart()
	//
	// defer m.Close()
	//
	// scoped := m.MustMock(
	// 	Get(URLPath("/test")).
	// 		Header("test", Should(Be(Equal("hello")))).
	// 		Query("filter", Is(Equal("all"))).
	// 		Reply(reply.
	// 			Created().
	// 			PlainText("hello world")))
	//
	// req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)
	// req.Header.Add("test", "hello")
	//
	// res, err := http.DefaultClient.Do(req)
	// assert.NoError(t, err)
	//
	// body, err := io.ReadAll(res.ParsedBody)
	//
	// assert.NoError(t, err)
	// assert.True(t, scoped.HasBeenCalled())
	// assert.Equal(t, 201, res.StatusCode)
	// assert.Equal(t, string(body), "hello world")
}

func TestMocha_NoReply(t *testing.T) {
	m := New()

	scoped, err := m.Mock(Get(URLPath("/test")))
	assert.Nil(t, scoped)
	assert.Error(t, err)
}

func TestMocha_NoMatchers(t *testing.T) {
	m := New()

	scoped, err := m.Mock(Request())
	assert.Nil(t, scoped)
	assert.Error(t, err)
}
