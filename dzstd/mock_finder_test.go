package dzstd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
	"github.com/vitorsalgado/mocha/v3/dzstd"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestMockMatches(t *testing.T) {
	params := &dzhttp.HTTPValueSelectorInput{}

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
			desc := &dzstd.Description{}
			pass, _ := dzstd.Match(params, desc, []*dzstd.Expectation[*dzhttp.HTTPValueSelectorInput]{{
				Matcher: StrictEqual(tc.value),
				ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
					return tc.selector
				},
			}})

			require.Equal(t, tc.expected, pass)
		})
	}

	t.Run("should return not matched and error when one of expectations returns error", func(t *testing.T) {
		// string
		desc := &dzstd.Description{}
		pass, _ := dzstd.Match(params, desc, []*dzstd.Expectation[*dzhttp.HTTPValueSelectorInput]{{
			Matcher: Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
				return "dev"
			},
		}})

		require.False(t, pass)
	})

	t.Run("should not pass when it panics", func(t *testing.T) {
		// string
		desc := &dzstd.Description{}
		pass, _ := dzstd.Match(params, desc, []*dzstd.Expectation[*dzhttp.HTTPValueSelectorInput]{{
			Matcher: Func(func(_ any) (bool, error) {
				panic("boom!")
			}),
			ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
				return "dev"
			},
		}})

		require.False(t, pass)
	})

	t.Run("should return the sum of the matchers Weight when it matches", func(t *testing.T) {
		// any
		desc := &dzstd.Description{}
		pass, weigth := dzstd.Match(params, desc, []*dzstd.Expectation[*dzhttp.HTTPValueSelectorInput]{
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
					return "test"
				},
				Weight: 1,
			},
			{
				Matcher: StrictEqual(10.0),
				ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
					return 10.0
				},
				Weight: 2,
			},
		})

		assert.True(t, pass)
		assert.Equal(t, 5, weigth)
	})

	t.Run("should return the sum of the matchers Weight when one of then doesn't match", func(t *testing.T) {
		// any
		desc := &dzstd.Description{}
		pass, weight := dzstd.Match(params, desc, []*dzstd.Expectation[*dzhttp.HTTPValueSelectorInput]{
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: StrictEqual("test"),
				ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
					return "dev"
				},
				Weight: 1,
			},
			{
				Matcher: StrictEqual(10.0),
				ValueSelector: func(r *dzhttp.HTTPValueSelectorInput) any {
					return 10.0
				},
				Weight: 2,
			},
		})

		assert.False(t, pass)
		assert.Equal(t, 4, weight)
	})
}
