package mocha

import (
	"fmt"
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
	params := &expect.RequestInfo{}

	cases := []struct {
		name     string
		value    any
		selector any
		expected bool
	}{
		{
			value:    "test",
			selector: "test",
			expected: true,
		},
		{
			value:    float64(10.0),
			selector: float64(10.0),
			expected: true,
		},
		{
			value:    true,
			selector: true,
			expected: true,
		},
		{
			value:    map[string]any{"key": "value"},
			selector: map[string]any{"key": "value"},
			expected: true,
		},
		{
			value:    "test",
			selector: "dev",
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := m.requestMatches(params, []Expectation{{
				Matcher: expect.ToEqual(tc.value),
				ValueSelector: func(r *expect.RequestInfo) any {
					return tc.selector
				},
			}})
			assert.Equal(t, tc.expected, res.OK)
		})
	}

	t.Run("should return not matched and error when one of expectations returns error", func(t *testing.T) {
		// string
		res := m.requestMatches(params, []Expectation{{
			Matcher: expect.Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			ValueSelector: func(r *expect.RequestInfo) any {
				return "dev"
			},
		}})
		assert.False(t, res.OK)
	})

	t.Run("should return the sum of the matchers weight when it matches", func(t *testing.T) {
		// any
		res := m.requestMatches(params, []Expectation{
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
		assert.True(t, res.OK)
		assert.Equal(t, 5, res.Weight)
	})

	t.Run("should return the sum of the matchers weight when one of then doesnt matches", func(t *testing.T) {
		// any
		res := m.requestMatches(params, []Expectation{
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
		assert.False(t, res.OK)
		assert.Equal(t, 4, res.Weight)
	})
}
