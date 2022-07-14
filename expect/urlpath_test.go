package expect

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlPath(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/test/hello")

	t.Run("should accept a non pointer", func(t *testing.T) {
		result, err := URLPath("/test/hello").Matches(*u, Args{})

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should accept a pointer", func(t *testing.T) {
		result, err := URLPath("/test/hello").Matches(u, Args{})

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should accept a string", func(t *testing.T) {
		su := "http://localhost:8080/test/hello"

		result, err := URLPath("/test/hello").Matches(su, Args{})

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should panic when providing a type that is not handled by URLPath", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = URLPath("/test/hello").Matches(10, Args{})
		})
	})

	t.Run("should return false when it doesnt match", func(t *testing.T) {
		result, err := URLPath("/test/bye").Matches(*u, Args{})

		assert.Nil(t, err)
		assert.False(t, result)
	})
}
