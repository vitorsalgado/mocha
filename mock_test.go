package mocha

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestRace(t *testing.T) {
	m := newMock()
	jobs := 10
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(index int) {
			m.Inc()
			m.Hits()
			wg.Done()
		}(i)

		m.Inc()
	}

	m.Hits()
	m.Inc()
	m.Inc()

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
		assert.False(t, m.HasBeenCalled())
		m.Inc()
		assert.True(t, m.HasBeenCalled())

		m.Dec()
		assert.False(t, m.HasBeenCalled())
	})
}

func TestMock_Matches(t *testing.T) {
	m := newMock()
	params := &valueSelectorInput{}

	testCases := []struct {
		name     string
		value    any
		selector any
		expected bool
	}{
		{value: "test", selector: "test", expected: true},
		{value: 10.0, selector: 10.0, expected: true},
		{value: true, selector: true, expected: true},
		{value: map[string]any{"key": "value"}, selector: map[string]any{"key": "value"}, expected: true},
		{value: "test", selector: "dev", expected: false},
		{value: make(chan struct{}, 1), selector: nil, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := m.matchExpectations(params, []*expectation{{
				Matcher: StrictEqual(tc.value),
				ValueSelector: func(r *valueSelectorInput) any {
					return tc.selector
				},
			}})
			assert.Equal(t, tc.expected, res.Pass)
		})
	}

	t.Run("should return not matched and error when one of expectations returns error", func(t *testing.T) {
		// string
		res := m.matchExpectations(params, []*expectation{{
			Matcher: Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			ValueSelector: func(r *valueSelectorInput) any {
				return "dev"
			},
		}})
		assert.False(t, res.Pass)
	})

	t.Run("should not pass when it panics", func(t *testing.T) {
		// string
		res := m.matchExpectations(params, []*expectation{{
			Matcher: Func(func(_ any) (bool, error) {
				panic("boom!")
			}),
			ValueSelector: func(r *valueSelectorInput) any {
				return "dev"
			},
		}})
		assert.False(t, res.Pass)
	})

	t.Run("should return the sum of the matchers weight when it matches", func(t *testing.T) {
		// any
		res := m.matchExpectations(params, []*expectation{
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *valueSelectorInput) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *valueSelectorInput) any {
					return "test"
				},
				Weight: 1,
			},
			{
				Matcher: StrictEqual(10.0),
				ValueSelector: func(r *valueSelectorInput) any {
					return 10.0
				},
				Weight: 2,
			},
		})
		assert.True(t, res.Pass)
		assert.Equal(t, 5, res.Weight)
	})

	t.Run("should return the sum of the matchers weight when one of then doesn't match", func(t *testing.T) {
		// any
		res := m.matchExpectations(params, []*expectation{
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *valueSelectorInput) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *valueSelectorInput) any {
					return "dev"
				},
				Weight: 1,
			},
			{
				Matcher: StrictEqual(10.0),
				ValueSelector: func(r *valueSelectorInput) any {
					return 10.0
				},
				Weight: 2,
			},
		})
		assert.False(t, res.Pass)
		assert.Equal(t, 4, res.Weight)
	})
}

func TestMock_Build(t *testing.T) {
	m := newMock()
	m.Inc()
	m.Disable()

	mm, err := m.Build()

	assert.NoError(t, err)
	assert.Equal(t, m, mm)
}
