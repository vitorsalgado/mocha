package core

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/core/mocks"
	"github.com/vitorsalgado/mocha/expect"
)

var fakeT = mocks.NewT()

func TestRace(t *testing.T) {
	m := NewMock()
	jobs := 10
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(index int) {
			if index%2 == 0 {
				time.Sleep(100 * time.Millisecond)
			}

			m.Hit()
			wg.Done()
		}(i)

		m.Hit()
	}

	m.Hit()
	m.Hit()

	wg.Wait()

	assert.Equal(t, (jobs*2)+2, m.Hits)
}

func TestMock(t *testing.T) {
	m := NewMock()

	t.Run("should init enabled", func(t *testing.T) {
		assert.True(t, m.Enabled)
	})

	t.Run("should disable mock when calling .Disable()", func(t *testing.T) {
		m.Disable()
		assert.False(t, m.Enabled)

		m.Enable()
		assert.True(t, m.Enabled)
	})

	t.Run("should return called when it was hit", func(t *testing.T) {
		assert.False(t, m.Called())
		m.Hit()
		assert.True(t, m.Called())

		m.Dec()
		assert.False(t, m.Called())
	})
}

func TestMock_Matches(t *testing.T) {
	m := NewMock()
	params := expect.Args{}

	t.Run("should return error when unable to cast the given expectation", func(t *testing.T) {
		res, err := m.Matches(params, []any{Expectation[int32]{}}, fakeT)
		assert.False(t, res.IsMatch)
		assert.Error(t, err)
	})

	t.Run("should match when generic type is known and matcher returns true without errors", func(t *testing.T) {
		// any
		res, err := m.Matches(params, []any{Expectation[any]{
			Matcher: expect.ToEqual[any]("test"),
			ValueSelector: func(r *expect.RequestInfo) any {
				return "test"
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// string
		res, err = m.Matches(params, []any{Expectation[string]{
			Matcher: expect.ToEqual("test"),
			ValueSelector: func(r *expect.RequestInfo) string {
				return "test"
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// float64
		res, err = m.Matches(params, []any{Expectation[float64]{
			Matcher: expect.ToEqual(10.0),
			ValueSelector: func(r *expect.RequestInfo) float64 {
				return 10.0
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// bool
		res, err = m.Matches(params, []any{Expectation[bool]{
			Matcher: expect.ToEqual(true),
			ValueSelector: func(r *expect.RequestInfo) bool {
				return true
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// map[string]any
		res, err = m.Matches(params, []any{Expectation[map[string]any]{
			Matcher: expect.ToEqual(map[string]any{"key": "value"}),
			ValueSelector: func(r *expect.RequestInfo) map[string]any {
				return map[string]any{"key": "value"}
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// map[string]any
		res, err = m.Matches(params, []any{Expectation[map[string][]string]{
			Matcher: expect.ToEqual(map[string][]string{"key": {"value1", "value2"}}),
			ValueSelector: func(r *expect.RequestInfo) map[string][]string {
				return map[string][]string{"key": {"value1", "value2"}}
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// []any]
		res, err = m.Matches(params, []any{Expectation[[]any]{
			Matcher: expect.ToEqual([]any{"test"}),
			ValueSelector: func(r *expect.RequestInfo) []any {
				return []any{"test"}
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// url.URL
		u, _ := url.Parse("http://localhost:8080")
		res, err = m.Matches(params, []any{Expectation[url.URL]{
			Matcher: expect.ToEqual(*u),
			ValueSelector: func(r *expect.RequestInfo) url.URL {
				return *u
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// url.Value
		res, err = m.Matches(params, []any{Expectation[url.Values]{
			Matcher: expect.ToEqual(url.Values{}),
			ValueSelector: func(r *expect.RequestInfo) url.Values {
				return url.Values{}
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// http.Request
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		res, err = m.Matches(params, []any{Expectation[*http.Request]{
			Matcher: expect.ToEqual(req),
			ValueSelector: func(r *expect.RequestInfo) *http.Request {
				return req
			},
		}}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)
	})

	t.Run("should return not matched result when one of expectations returns false", func(t *testing.T) {
		// string
		res, err := m.Matches(params, []any{Expectation[string]{
			Matcher: expect.ToEqual("test"),
			ValueSelector: func(r *expect.RequestInfo) string {
				return "dev"
			},
		}}, fakeT)
		assert.False(t, res.IsMatch)
		assert.Nil(t, err)
	})

	t.Run("should return not matched and error when one of expectations returns error", func(t *testing.T) {
		// string
		res, err := m.Matches(params, []any{Expectation[string]{
			Matcher: expect.Func(func(_ string, p expect.Args) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			ValueSelector: func(r *expect.RequestInfo) string {
				return "dev"
			},
		}}, fakeT)
		assert.False(t, res.IsMatch)
		assert.NotNil(t, err)
	})

	t.Run("should return the sum of the matchers weight when it matches", func(t *testing.T) {
		// any
		res, err := m.Matches(params, []any{
			Expectation[any]{
				Matcher: expect.ToEqualAny("test"),
				ValueSelector: func(r *expect.RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			Expectation[string]{
				Matcher: expect.ToEqual("test"),
				ValueSelector: func(r *expect.RequestInfo) string {
					return "test"
				},
				Weight: 1,
			},
			Expectation[float64]{
				Matcher: expect.ToEqual(10.0),
				ValueSelector: func(r *expect.RequestInfo) float64 {
					return 10.0
				},
				Weight: 2,
			},
		}, fakeT)
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)
		assert.Equal(t, 5, res.Weight)
	})

	t.Run("should return the sum of the matchers weight when one of then doesnt matches", func(t *testing.T) {
		// any
		res, err := m.Matches(params, []any{
			Expectation[any]{
				Matcher: expect.ToEqualAny("test"),
				ValueSelector: func(r *expect.RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			Expectation[string]{
				Matcher: expect.ToEqual("test"),
				ValueSelector: func(r *expect.RequestInfo) string {
					return "dev"
				},
				Weight: 1,
			},
			Expectation[float64]{
				Matcher: expect.ToEqual(10.0),
				ValueSelector: func(r *expect.RequestInfo) float64 {
					return 10.0
				},
				Weight: 2,
			},
		}, fakeT)
		assert.False(t, res.IsMatch)
		assert.Nil(t, err)
		assert.Equal(t, 5, res.Weight)
	})
}
