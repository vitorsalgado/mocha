package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEither(t *testing.T) {
	t.Parallel()

	t.Run("should return true when only left matcher evaluates to true", func(t *testing.T) {
		result, err := Either(EqualTo("test"), Contains("qa"))("test", Args{})
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when only right matcher evaluates to true", func(t *testing.T) {
		result, err := Either(EqualTo("qa"), Contains("tes"))("test", Args{})
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when both matchers evaluates to true", func(t *testing.T) {
		result, err := Either(EqualTo("test"), Contains("te"))("test", Args{})
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when both evaluates to false", func(t *testing.T) {
		result, err := Either(EqualTo("dev"), Contains("qa"))("test", Args{})
		assert.Nil(t, err)
		assert.False(t, result)
	})
}
