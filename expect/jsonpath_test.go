package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _json = map[string]any{
	"name": "someone",
	"age":  34,
	"address": map[string]any{
		"street": "very nice place",
	},
	"job": nil,
}

func TestJSONPathMatcher(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		matcher  Matcher
		expected bool
	}{
		{
			name:     "should read match text field on object root",
			path:     "name",
			matcher:  ToEqual("someone"),
			expected: true,
		},
		{
			name:     "should match numeric field value",
			path:     "age",
			matcher:  ToEqual(34),
			expected: true,
		},
		{
			name:     "should match nested object field",
			path:     "address.street",
			matcher:  ToEqual("very nice place"),
			expected: true,
		},
		{
			name:     "should match nil when field is present",
			path:     "job",
			matcher:  ToEqual(nil),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := JSONPath(tc.path, tc.matcher).Match(_json)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, res.OK)
		})
	}
}

func TestJSONPathMatcher_Match_Errors(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		matcher Matcher
	}{
		{name: "invalid path", path: "312nj.,", matcher: ToEqual("anything")},
		{name: "field not present", path: "life", matcher: ToEqual("anything")},
		{name: "deep field not present", path: "path.to.another.field", matcher: ToEqual("any")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := JSONPath(tc.path, tc.matcher).Match(_json)
			assert.NotNil(t, err)
			assert.False(t, res.OK)
		})
	}
}
