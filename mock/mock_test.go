package mock

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/matcher"
)

func TestRace(t *testing.T) {
	m := New()
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
	m := New()

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
	m := New()
	params := matcher.Args{}

	t.Run("should return error when unable to cast the given expectation", func(t *testing.T) {
		res, err := m.Matches(params, []any{Expectation[int32]{}})
		assert.False(t, res.IsMatch)
		assert.Error(t, err)
	})

	t.Run("should match when generic type is known and matcher returns true without errors", func(t *testing.T) {
		// any
		res, err := m.Matches(params, []any{Expectation[any]{
			Matcher: matcher.EqualAny("test"),
			ValuePicker: func(r *matcher.RequestInfo) any {
				return "test"
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// string
		res, err = m.Matches(params, []any{Expectation[string]{
			Matcher: matcher.EqualTo("test"),
			ValuePicker: func(r *matcher.RequestInfo) string {
				return "test"
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// float64
		res, err = m.Matches(params, []any{Expectation[float64]{
			Matcher: matcher.EqualTo(10.0),
			ValuePicker: func(r *matcher.RequestInfo) float64 {
				return 10.0
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// bool
		res, err = m.Matches(params, []any{Expectation[bool]{
			Matcher: matcher.EqualTo(true),
			ValuePicker: func(r *matcher.RequestInfo) bool {
				return true
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// map[string]any
		res, err = m.Matches(params, []any{Expectation[map[string]any]{
			Matcher: matcher.EqualTo(map[string]any{"key": "value"}),
			ValuePicker: func(r *matcher.RequestInfo) map[string]any {
				return map[string]any{"key": "value"}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// map[string]any
		res, err = m.Matches(params, []any{Expectation[map[string][]string]{
			Matcher: matcher.EqualTo(map[string][]string{"key": {"value1", "value2"}}),
			ValuePicker: func(r *matcher.RequestInfo) map[string][]string {
				return map[string][]string{"key": {"value1", "value2"}}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// []any]
		res, err = m.Matches(params, []any{Expectation[[]any]{
			Matcher: matcher.EqualTo([]any{"test"}),
			ValuePicker: func(r *matcher.RequestInfo) []any {
				return []any{"test"}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// url.URL
		u, _ := url.Parse("http://localhost:8080")
		res, err = m.Matches(params, []any{Expectation[url.URL]{
			Matcher: matcher.EqualTo(*u),
			ValuePicker: func(r *matcher.RequestInfo) url.URL {
				return *u
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// url.Value
		res, err = m.Matches(params, []any{Expectation[url.Values]{
			Matcher: matcher.EqualTo(url.Values{}),
			ValuePicker: func(r *matcher.RequestInfo) url.Values {
				return url.Values{}
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)

		// http.Request
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		res, err = m.Matches(params, []any{Expectation[*http.Request]{
			Matcher: matcher.EqualTo(req),
			ValuePicker: func(r *matcher.RequestInfo) *http.Request {
				return req
			},
		}})
		assert.True(t, res.IsMatch)
		assert.Nil(t, err)
	})

	t.Run("should return not matched result when one of expectations returns false", func(t *testing.T) {
		// string
		res, err := m.Matches(params, []any{Expectation[string]{
			Matcher: matcher.EqualTo("test"),
			ValuePicker: func(r *matcher.RequestInfo) string {
				return "dev"
			},
		}})
		assert.False(t, res.IsMatch)
		assert.Nil(t, err)
	})

	t.Run("should return not matched and error when one of expectations returns error", func(t *testing.T) {
		// string
		res, err := m.Matches(params, []any{Expectation[string]{
			Matcher: func(_ string, p matcher.Args) (bool, error) {
				return false, fmt.Errorf("fail")
			},
			ValuePicker: func(r *matcher.RequestInfo) string {
				return "dev"
			},
		}})
		assert.False(t, res.IsMatch)
		assert.NotNil(t, err)
	})

	t.Run("should return the sum of the matchers weight when it matches", func(t *testing.T) {
		// any
		res, err := m.Matches(params, []any{
			Expectation[any]{
				Matcher: matcher.EqualAny("test"),
				ValuePicker: func(r *matcher.RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			Expectation[string]{
				Matcher: matcher.EqualTo("test"),
				ValuePicker: func(r *matcher.RequestInfo) string {
					return "test"
				},
				Weight: 1,
			},
			Expectation[float64]{
				Matcher: matcher.EqualTo(10.0),
				ValuePicker: func(r *matcher.RequestInfo) float64 {
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
		res, err := m.Matches(params, []any{
			Expectation[any]{
				Matcher: matcher.EqualAny("test"),
				ValuePicker: func(r *matcher.RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			Expectation[string]{
				Matcher: matcher.EqualTo("test"),
				ValuePicker: func(r *matcher.RequestInfo) string {
					return "dev"
				},
				Weight: 1,
			},
			Expectation[float64]{
				Matcher: matcher.EqualTo(10.0),
				ValuePicker: func(r *matcher.RequestInfo) float64 {
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
