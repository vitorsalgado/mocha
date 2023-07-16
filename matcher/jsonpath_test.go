package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var _json = map[string]any{
	"name": "someone",
	"age":  34,
	"address": map[string]any{
		"street": "very nice place",
	},
	"job": nil,
}

func TestJSONPath(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		matcher  Matcher
		expected bool
	}{
		{
			name:     "should read match text field on object root",
			path:     "name",
			matcher:  StrictEqual("someone"),
			expected: true,
		},
		{
			name:     "should match numeric field value",
			path:     "age",
			matcher:  StrictEqual(34),
			expected: true,
		},
		{
			name:     "should match nested object field",
			path:     "address.street",
			matcher:  StrictEqual("very nice place"),
			expected: true,
		},
		{
			name:     "should match nil when field is present",
			path:     "job",
			matcher:  StrictEqual(nil),
			expected: true,
		},
		{
			name:     "not present matchers for absent field",
			path:     "address.city",
			matcher:  Not(Present()),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := JSONPath(tc.path, tc.matcher).Match(_json)
			require.NoError(t, err)
			require.Equal(t, tc.expected, res.Pass)
		})
	}
}

func TestJSONPathMatcherMatchErrors(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		matcher Matcher
	}{
		{name: "field not present", path: "life", matcher: StrictEqual("anything")},
		{name: "deep field not present", path: "path.to.another.field", matcher: StrictEqual("any")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := JSONPath(tc.path, tc.matcher).Match(_json)
			require.NoError(t, err)
			require.False(t, res.Pass)
		})
	}
}

func TestJSONPathMatcherInvalidPaths(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		matcher Matcher
	}{
		{name: "invalid path", path: "312nj.,", matcher: StrictEqual("anything")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Panics(t, func() {
				_, _ = JSONPath(tc.path, tc.matcher).Match(_json)
			})
		})
	}
}

func TestJSONPathNew(t *testing.T) {
	require.Panics(t, func() {
		JSONPath(".", Eq(""))
	})

	require.Panics(t, func() {
		Field(".", Eq(""))
	})
}
