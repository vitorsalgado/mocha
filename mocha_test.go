package mocha

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	mok "github.com/stretchr/testify/mock"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/reply"
	"github.com/vitorsalgado/mocha/to"
)

type (
	TestModel struct {
		Name string `json:"name"`
		OK   bool   `json:"ok"`
	}

	FakeT struct{ mok.Mock }
)

func (m *FakeT) Cleanup(_ func()) {
	m.Called()
}

func (m *FakeT) Helper() {
	m.Called()
}

func (m *FakeT) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func TestMocha(t *testing.T) {
	t.Run("should mock request", func(t *testing.T) {
		m := ForTest(t)
		m.Start()

		scoped := m.Mock(
			Get(to.HaveURLPath("/test")).
				Header("test", to.Equal("hello")).
				Query("filter", to.Equal("all")).
				Reply(reply.
					Created().
					BodyString("hello world")))

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

func TestPostJSON(t *testing.T) {
	m := ForTest(t)
	m.Start()

	scoped := m.Mock(Post(to.HaveURLPath("/test")).
		Header("test", to.Equal("hello")).
		Body(
			to.JSONPath("name", to.Equal("dev")), to.JSONPath("ok", to.Equal(true))).
		Reply(reply.OK()))

	req := testutil.PostJSON(m.Server.URL+"/test", &TestModel{Name: "dev", OK: true})
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

	scope := m.Mock(Get(to.HaveURLPath("/test")).
		Matches(to.Fn(func(v any, params to.Args) (bool, error) {
			p, _ := params.Params.Get(key)
			return p.(string) == expected, nil
		})).
		Reply(reply.Accepted()))

	req := testutil.Get(fmt.Sprintf("%s/test", m.Server.URL))
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	scope.MustBeDone()
	assert.Equal(t, http.StatusAccepted, res.StatusCode)
}

func TestResponseMapper(t *testing.T) {
	m := ForTest(t)
	m.Start()

	scoped := m.Mock(Get(to.HaveURLPath("/test")).
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

	scoped.MustBeDone()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "dev", res.Header.Get("x-test"))
}

func TestDelay(t *testing.T) {
	m := ForTest(t)
	m.Start()

	start := time.Now()
	delay := time.Duration(1250) * time.Millisecond

	scoped := m.Mock(Get(to.HaveURLPath("/test")).
		Reply(reply.
			OK().
			Delay(delay)))

	req := testutil.Get(fmt.Sprintf("%s/test", m.Server.URL))
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	elapsed := time.Since(start)

	scoped.MustBeDone()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.GreaterOrEqual(t, elapsed, delay)
}

func TestPostExpectations(t *testing.T) {
	m := ForTest(t)
	m.Start()

	scoped := m.Mock(
		NewBuilder().
			MatchAfter(to.Repeat(2)).
			Method("GET").
			URL(to.HaveURLPath("/test")).
			Reply(reply.
				OK()))

	testutil.Get(fmt.Sprintf("%s/other", m.Server.URL)).Do()
	testutil.Get(fmt.Sprintf("%s/other", m.Server.URL)).Do()

	res, _ := testutil.Get(fmt.Sprintf("%s/other", m.Server.URL)).Do()
	assert.Equal(t, res.StatusCode, http.StatusTeapot)

	res, _ = testutil.Get(fmt.Sprintf("%s/test", m.Server.URL)).Do()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, _ = testutil.Get(fmt.Sprintf("%s/test", m.Server.URL)).Do()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	res, _ = testutil.Get(fmt.Sprintf("%s/test", m.Server.URL)).Do()
	assert.Equal(t, res.StatusCode, http.StatusTeapot)

	scoped.MustBeDone()
}

func TestErrors(t *testing.T) {
	fake := &FakeT{}

	fake.On("Cleanup", mok.Anything).Return()
	fake.On("Helper").Return()
	fake.On("Errorf", mok.AnythingOfType("string"), mok.Anything).Return()

	m := ForTest(fake)
	m.Start()

	defer m.Close()

	t.Run("should log errors on reply", func(t *testing.T) {
		scoped := m.Mock(Get(to.HaveURLPath("/test1")).
			ReplyFunction(func(r *http.Request, m *mock.Mock, p params.Params) (*mock.Response, error) {
				return nil, fmt.Errorf("failed to build a response")
			}))

		res, err := testutil.Get(fmt.Sprintf("%s/test1", m.Server.URL)).Do()

		assert.Nil(t, err)
		assert.True(t, scoped.IsDone())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		fake.AssertExpectations(t)
	})

	t.Run("should log errors from matchers", func(t *testing.T) {
		scoped := m.Mock(Get(to.HaveURLPath("/test2")).
			Header("test", to.Fn(
				func(_ string, _ to.Args) (bool, error) {
					return false, fmt.Errorf("failed")
				})))

		res, err := testutil.Get(fmt.Sprintf("%s/test2", m.Server.URL)).Do()

		assert.Nil(t, err)
		assert.False(t, scoped.IsDone())
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		fake.AssertExpectations(t)
	})
}
