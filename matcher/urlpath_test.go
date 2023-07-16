package matcher

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLPath(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/test/hello")

	testCases := []struct {
		name     string
		path     string
		value    any
		expected bool
	}{
		{"should accept a pointer", "/test/hello", u, true},
		{"should accept a string", "/test/hello", *u, true},
		{"should return false when it doesn't match", "/test/bye", u.String(), false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := URLPath(tt.path).Match(u)

			assert.Nil(t, err)
			assert.Equal(t, tt.expected, result.Pass)

			result, err = URLPathMatch(Contain(tt.path)).Match(u)

			assert.Nil(t, err)
			assert.Equal(t, tt.expected, result.Pass)
		})
	}

	t.Run("should return error when providing a type that is not handled by URLPath", func(t *testing.T) {
		res, err := URLPath("/test/hello").Match(10)
		require.Error(t, err)
		require.False(t, res.Pass)
	})
}
