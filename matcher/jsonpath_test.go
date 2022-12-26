package matcher

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
			matcher:  Equal("someone"),
			expected: true,
		},
		{
			name:     "should match numeric field value",
			path:     "age",
			matcher:  Equal(34),
			expected: true,
		},
		{
			name:     "should match nested object field",
			path:     "address.street",
			matcher:  Equal("very nice place"),
			expected: true,
		},
		{
			name:     "should match nil when field is present",
			path:     "job",
			matcher:  Equal(nil),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := JSONPath(tc.path, tc.matcher).Match(_json)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, res.Pass)
		})
	}
}

func TestJSONPathMatcher_Match_Errors(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		matcher Matcher
	}{
		{name: "invalid path", path: "312nj.,", matcher: Equal("anything")},
		{name: "field not present", path: "life", matcher: Equal("anything")},
		{name: "deep field not present", path: "path.to.another.field", matcher: Equal("any")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := JSONPath(tc.path, tc.matcher).Match(_json)
			assert.Nil(t, res)
			assert.Error(t, err)
		})
	}
}
