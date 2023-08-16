package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlob(t *testing.T) {
	testCases := []struct {
		pattern        string
		value          any
		expectedResult bool
	}{
		{"*", "/test", true},
		{"**", "*", true},
		{"**", "**", true},
		{"/test/*dev", "/test/*dev", true},
		{"/*", "/", true},
		{"/*", "/test", true},
		{"/test", "/test", true},
		{"/test", "/TEST", false},
		{"/test", "/dev", false},
		{"/test?q=test&f=all", "/test?q=test&f=all", true},
		{"/test?q=test&f=none", "/test?q=test&f=all", false},
		{"/test*", "/test*", true},
		{"/test/*", "/test/", true},
		{"/test*", "/test", true},
		{"/test*", "/test?q=test&f=all", true},
		{"/test/**", "/test/test2/test3", true},
		{"/test/**", "/test/TEST2/TEST3", true},
		{"/test/**/world", "/test/hello/world", true},
		{"https://www.example.org/test/*", "https://www.example.org/test/hello/world", true},
		{"https://www.example.org/test/**", "https://www.example.org/test/hello/world", true},
		{"https://www.example.org/test/**/world", "https://www.example.org/test/hello/world", true},
		{"**/test/**", "https://www.example.org/test/hello/world", true},
		{"**/test/*", "https://www.example.org/test/hello/world", true},
		{"/test/**", "https://www.example.org/test/hello/world", false},
		{"/test/*", "https://www.example.org/test/hello/world", false},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern, func(t *testing.T) {
			result, err := GlobMatch(tc.pattern).Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, result.Pass)
		})
	}
}
