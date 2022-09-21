package mocha

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/expect"
)

func TestRace(t *testing.T) {
	m := newMock()
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

	assert.Equal(t, (jobs*2)+2, m.hits)
}

func TestMock(t *testing.T) {
	m := newMock()

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
	m := newMock()
	params := expect.Args{}

	t.Run("should match when generic type is known and matcher returns true without errors", func(t *testing.T) {
		// any
		res, err := m.matches(params, []Expectation{{
			Matcher: expect.ToEqual("test"),
			ValueSelector: func(r *expect.RequestInfo) any {
				return "test"
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// string
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual("test"),
			ValueSelector: func(r *expect.RequestInfo) any {
				return "test"
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// float64
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual(10.0),
			ValueSelector: func(r *expect.RequestInfo) any {
				return 10.0
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// bool
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual(true),
			ValueSelector: func(r *expect.RequestInfo) any {
				return true
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// map[string]any
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual(map[string]any{"key": "value"}),
			ValueSelector: func(r *expect.RequestInfo) any {
				return map[string]any{"key": "value"}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// map[string]any
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual(map[string][]string{"key": {"value1", "value2"}}),
			ValueSelector: func(r *expect.RequestInfo) any {
				return map[string][]string{"key": {"value1", "value2"}}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// []any]
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual([]any{"test"}),
			ValueSelector: func(r *expect.RequestInfo) any {
				return []any{"test"}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// url.URL
		u, _ := url.Parse("http://localhost:8080")
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual(*u),
			ValueSelector: func(r *expect.RequestInfo) any {
				return *u
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// url.Value
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual(url.Values{}),
			ValueSelector: func(r *expect.RequestInfo) any {
				return url.Values{}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// http.Request
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		res, err = m.matches(params, []Expectation{{
			Matcher: expect.ToEqual(req),
			ValueSelector: func(r *expect.RequestInfo) any {
				return req
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)
	})

	t.Run("should return not matched result when one of expectations returns false", func(t *testing.T) {
		// string
		res, err := m.matches(params, []Expectation{{
			Matcher: expect.ToEqual("test"),
			ValueSelector: func(r *expect.RequestInfo) any {
				return "dev"
			},
		}})
		assert.False(t, res.IsMatch)
		assert.Nil(t, err)
	})

	t.Run("should return not matched and error when one of expectations returns error", func(t *testing.T) {
		// string
		res, err := m.matches(params, []Expectation{{
			Matcher: expect.Func(func(_ any, p expect.Args) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			ValueSelector: func(r *expect.RequestInfo) any {
				return "dev"
			},
		}})
		assert.False(t, res.IsMatch)
		assert.NotNil(t, err)
	})

	t.Run("should return the sum of the matchers weight when it matches", func(t *testing.T) {
		// any
		res, err := m.matches(params, []Expectation{
			{
				Matcher: expect.ToEqual("test"),
				ValueSelector: func(r *expect.RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: expect.ToEqual("test"),
				ValueSelector: func(r *expect.RequestInfo) any {
					return "test"
				},
				Weight: 1,
			},
			{
				Matcher: expect.ToEqual(10.0),
				ValueSelector: func(r *expect.RequestInfo) any {
					return 10.0
				},
				Weight: 2,
			},
		})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)
		assert.Equal(t, 5, res.Weight)
	})

	t.Run("should return the sum of the matchers weight when one of then doesnt matches", func(t *testing.T) {
		// any
		res, err := m.matches(params, []Expectation{
			{
				Matcher: expect.ToEqual("test"),
				ValueSelector: func(r *expect.RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: expect.ToEqual("test"),
				ValueSelector: func(r *expect.RequestInfo) any {
					return "dev"
				},
				Weight: 1,
			},
			{
				Matcher: expect.ToEqual(10.0),
				ValueSelector: func(r *expect.RequestInfo) any {
					return 10.0
				},
				Weight: 2,
			},
		})
		assert.False(t, res.IsMatch)
		assert.Nil(t, err)
		assert.Equal(t, 5, res.Weight)
	})
}
