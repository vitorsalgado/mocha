package mocha

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	. "github.com/vitorsalgado/mocha/v3/matcher"
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
			Header("test", StrictEqual("hello")).
			Query("filter", StrictEqual("all")).
			Reply(Created().
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

func TestResponseMapper_ModifyingResponse(t *testing.T) {
	const k = "key"
	const v = "test-ok"

	m := New()
	_ = m.Parameters().Set(context.Background(), k, v)
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(Get(URLPath("/test")).
		Reply(OK()).
		Map(func(rv *RequestValues, r *Stub) error {
			val, _, _ := rv.App.Parameters().Get(rv.RawRequest.Context(), k)

			r.Header.Add("x-param-key", val.(string))
			r.Header.Add("x-test", rv.RawRequest.Header.Get("x-param"))

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
	assert.Equal(t, "test-ok", res.Header.Get("x-param-key"))
}

func TestResponseDelay(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	start := time.Now()
	delay := 250 * time.Millisecond

	scoped := m.MustMock(Get(URLPath("/test")).
		Delay(delay).
		Reply(OK()))

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
			Reply(OK()))

	res, err := testutil.Get(fmt.Sprintf("%s/test2", m.URL())).Do()

	assert.NoError(t, err)
	assert.False(t, scoped.HasBeenCalled())
	assert.Equal(t, StatusNoMatch, res.StatusCode)
}

func TestMocha_Assertions(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	fakeT := newFakeT()

	scoped := m.MustMock(
		Get(URLPath("/test-ok")).
			Reply(OK()))

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
			Reply(OK()),
		Get(URLPath("/test-2")).
			Reply(OK()))

	res, err := testutil.Get(m.URL() + "/test-1").Do()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// disable all mocks
	// should return tea pot for all
	m.Disable()

	res, err = testutil.Get(m.URL() + "/test-1").Do()

	assert.NoError(t, err)
	assert.Equal(t, StatusNoMatch, res.StatusCode)

	res, err = testutil.Get(m.URL() + "/test-2").Do()

	assert.NoError(t, err)
	assert.Equal(t, StatusNoMatch, res.StatusCode)

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
			Reply(OK()))

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
			Reply(Created().
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

func TestSchemeMatching(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test")).
			Scheme("http").
			Reply(OK()))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	assert.True(t, scoped.HasBeenCalled())
	assert.Equal(t, http.StatusOK, res.StatusCode)
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

func TestMocha_Concurrency(t *testing.T) {
	m := New()
	jobs := 10
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(index int) {
			m.MustMock(Getf("/test--" + strconv.FormatInt(int64(index), 10)).Reply(OK()))
			m.MustMock(Getf("/concurrency--" + strconv.FormatInt(int64(index), 10)).Reply(Accepted()))

			m.Hits()
			m.Disable()
			m.Enable()

			m.Clean()

			m.MustMock(Getf("/test--" + strconv.FormatInt(int64(index), 10)).Reply(OK()))
			m.MustMock(Getf("/concurrency--" + strconv.FormatInt(int64(index), 10)).Reply(Accepted()))

			wg.Done()
		}(i)

		m.MustMock(Getf("/test-outside--" + strconv.FormatInt(int64(i), 10)).Reply(Created()))
		m.Clean()

		m.Hits()
		m.Disable()
		m.Enable()
	}

	m.Hits()
	m.MustMock(Getf("/test-final").Reply(BadRequest()))
	m.MustStart()
	m.Close()

	wg.Wait()
}

func TestMocha_Concurrent_Requests(t *testing.T) {
	jobs := 20
	wg := sync.WaitGroup{}
	httpClient := &http.Client{}

	m := New()
	m.MustStart()

	defer m.Close()

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(index int) {
			num := strconv.FormatInt(int64(index), 10)

			scope1 := m.MustMock(Getf("/test--" + num).Reply(OK()))
			scope2 := m.MustMock(Getf("/concurrency--" + num).Reply(Accepted()))
			scope3 := m.MustMock(Getf("/seq--" + num).Reply(Seq(OK(), BadRequest())))
			scope4 := m.MustMock(Getf("/rand--" + num).Reply(Rand(Accepted(), InternalServerError())))

			res, err := httpClient.Get(m.URL() + "/test--" + num)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			m.MustMock(Getf("/test-after--" + num).Reply(BadRequest()))

			res, err = httpClient.Get(m.URL() + "/concurrency--" + num)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusAccepted, res.StatusCode)

			m.MustMock(Getf("/concurrency-after--" + num).Reply(InternalServerError()))

			res, err = httpClient.Get(m.URL() + "/seq--" + num)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			res, err = httpClient.Get(m.URL() + "/seq--" + num)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)

			for n := 0; n < 50; n++ {
				res, err = httpClient.Get(m.URL() + "/rand--" + num)
				assert.NoError(t, err)
				assert.True(t, res.StatusCode == http.StatusAccepted || res.StatusCode == http.StatusInternalServerError)
			}

			scope1.AssertCalled(t)
			scope1.AssertNumberOfCalls(t, 1)
			scope2.AssertCalled(t)
			scope2.AssertNumberOfCalls(t, 1)
			scope3.AssertCalled(t)
			scope3.AssertNumberOfCalls(t, 2)
			scope4.AssertCalled(t)
			scope4.AssertNumberOfCalls(t, 50)

			m.MustMock(Getf("/test-final--" + num).Reply(PartialContent()))

			wg.Done()
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 54*jobs, m.Hits())
}
