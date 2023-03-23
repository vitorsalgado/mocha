package mocha

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/test/testmock"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(
		m,
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
	)
}

func TestMocha(t *testing.T) {
	m := NewAPI()
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
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.True(t, scoped.HasBeenCalled())
	require.Equal(t, 201, res.StatusCode)
	require.Equal(t, string(body), "hello world")
}

func TestResponseMapperModifyingResponse(t *testing.T) {
	const k = "key"
	const v = "test-ok"

	m := NewAPI()
	_ = m.Parameters().Set(k, v)
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(Get(URLPath("/test")).
		Reply(OK()).
		Map(func(rv *RequestValues, r *Stub) error {
			val, _ := rv.App.Parameters().Get(k)

			r.Header.Add("x-param-key", val.(string))
			r.Header.Add("x-test", rv.RawRequest.Header.Get("x-param"))

			return nil
		}))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/test", m.URL()), nil)
	req.Header.Add("x-param", "dev")

	res, err := http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "dev", res.Header.Get("x-test"))
	require.Equal(t, "test-ok", res.Header.Get("x-param-key"))
	require.True(t, scoped.AssertCalled(t))
}

func TestErrors(t *testing.T) {
	m := NewAPI()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test2")).
			Header("test", Func(
				func(_ any) (bool, error) {
					return false, fmt.Errorf("failed")
				})).
			Reply(OK()))

	res, err := http.Get(fmt.Sprintf("%s/test2", m.URL()))

	require.NoError(t, err)
	require.False(t, scoped.HasBeenCalled())
	require.Equal(t, StatusNoMatch, res.StatusCode)
}

func TestMochaAssertions(t *testing.T) {
	m := NewAPI()
	m.MustStart()

	defer m.Close()

	ft := testmock.NewFakeT()

	scoped := m.MustMock(
		Get(URLPath("/test-ok")).
			Reply(OK()))

	require.Equal(t, 0, scoped.Hits())
	require.False(t, m.AssertCalled(ft))
	require.True(t, m.AssertNotCalled(ft))
	require.True(t, m.AssertNumberOfCalls(ft, 0))
	require.Equal(t, 0, m.Hits())

	res, err := http.Get(m.URL() + "/test-ok")

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 1, scoped.Hits())
	require.True(t, m.AssertCalled(ft))
	require.False(t, m.AssertNotCalled(ft))
	require.True(t, m.AssertNumberOfCalls(ft, 1))
	require.Equal(t, 1, m.Hits())
}

func TestMochaEnableDisable(t *testing.T) {
	m := NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(
		Get(URLPath("/test-1")).
			Reply(OK()),
		Get(URLPath("/test-2")).
			Reply(OK()))

	res, err := http.Get(m.URL() + "/test-1")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	// disable all mocks
	// should return tea pot for all
	m.Disable()

	res, err = http.Get(m.URL() + "/test-1")
	require.NoError(t, err)
	require.Equal(t, StatusNoMatch, res.StatusCode)

	res, err = http.Get(m.URL() + "/test-2")
	require.NoError(t, err)
	require.Equal(t, StatusNoMatch, res.StatusCode)

	// re-enable mocks again
	m.Enable()

	res, err = http.Get(m.URL() + "/test-1")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get(m.URL() + "/test-2")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestMochaSilently(t *testing.T) {
	m := NewAPI(Setup().LogVerbosity(LogBasic).LogLevel(LogLevelNone))
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test")).
			Reply(Created().
				PlainText("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.True(t, scoped.HasBeenCalled())
	require.Equal(t, 201, res.StatusCode)
	require.Equal(t, string(body), "hello world")
}

func TestSchemeMatching(t *testing.T) {
	m := NewAPI()
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

	require.True(t, scoped.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestMochaNoReply(t *testing.T) {
	m := NewAPI()

	scoped, err := m.Mock(Get(URLPath("/test")))
	require.Nil(t, scoped)
	require.Error(t, err)
}

func TestMochaNoMatchers(t *testing.T) {
	m := NewAPI()

	scoped, err := m.Mock(Request())
	require.Nil(t, scoped)
	require.Error(t, err)
}

func TestMocha_RequestMatches(t *testing.T) {
	m := NewAPI()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		Get(URLPath("/test")).
			RequestMatches(func(r *http.Request) (bool, error) {
				if r.Method == http.MethodGet {
					return true, nil
				}

				return false, nil
			}).
			Reply(OK()))

	httpClient := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := httpClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	require.True(t, scoped.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)

	req, _ = http.NewRequest(http.MethodPost, m.URL()+"/test", nil)
	res, err = httpClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	scoped.AssertNumberOfCalls(t, 1)
	require.Equal(t, StatusNoMatch, res.StatusCode)
}

func TestMochaConcurrency(t *testing.T) {
	m := NewAPI()
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

func TestMochaConcurrentRequests(t *testing.T) {
	jobs := 20
	wg := sync.WaitGroup{}
	httpClient := &http.Client{}

	m := NewAPI()
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
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)

			m.MustMock(Getf("/test-after--" + num).Reply(BadRequest()))

			res, err = httpClient.Get(m.URL() + "/concurrency--" + num)
			require.NoError(t, err)
			require.Equal(t, http.StatusAccepted, res.StatusCode)

			m.MustMock(Getf("/concurrency-after--" + num).Reply(InternalServerError()))

			res, err = httpClient.Get(m.URL() + "/seq--" + num)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)

			res, err = httpClient.Get(m.URL() + "/seq--" + num)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, res.StatusCode)

			for n := 0; n < 50; n++ {
				res, err = httpClient.Get(m.URL() + "/rand--" + num)
				require.NoError(t, err)
				require.True(t, res.StatusCode == http.StatusAccepted || res.StatusCode == http.StatusInternalServerError)
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

	require.Equal(t, 54*jobs, m.Hits())
}

func TestSettingOnlyPort(t *testing.T) {
	randomPort := func() int {
		l, err := net.Listen("tcp", ":0")
		require.NoError(t, err)
		require.NoError(t, l.Close())

		return l.Addr().(*net.TCPAddr).Port
	}

	port := randomPort()
	m := NewAPIWithT(t, Setup().Port(port))
	m.MustStart()
	m.MustMock(Getf("/test").Reply(OK()))

	hc := &http.Client{}
	res, err := hc.Get(m.URL("/test"))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Contains(t, m.URL(), strconv.FormatInt(int64(port), 10))
}

func TestMocha_Server(t *testing.T) {
	m := NewAPI()
	srv := m.Server()

	require.NotNil(t, srv)
	require.IsType(t, &httpTestServer{}, srv)
	require.IsType(t, &httptest.Server{}, srv.S())
}
