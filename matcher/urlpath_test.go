package matcher

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLPath(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/test/hello")

	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{"should accept a pointer", "/test/hello", true},
		{"should accept a string", "/test/hello", true},
		{"should return false when it doesnt match", "/test/bye", false},
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

	t.Run("should panic when providing a type that is not handled by URLPath", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = URLPath("/test/hello").Match(10)
		})
	})
}
