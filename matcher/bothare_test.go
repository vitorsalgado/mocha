package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBothAre(t *testing.T) {
	t.Parallel()

	t.Run("should return false when only left matcher evaluates to true", func(t *testing.T) {
		result, err := BothAre(EqualTo("test")).And(Contains("qa"))("test", Args{})
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when only right matcher evaluates to true", func(t *testing.T) {
		result, err := BothAre(EqualTo("qa")).And(Contains("tes"))("test", Args{})
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when both matchers evaluates to true", func(t *testing.T) {
		result, err := BothAre(EqualTo("test")).And(Contains("te"))("test", Args{})
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when only left matcher evaluates to true", func(t *testing.T) {
		result, err := BothAre(EqualTo("test")).And(Contains("qa"))("test", Args{})
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when matchers throws errors", func(t *testing.T) {
		result, err := BothAre(
			func(_ string, _ Args) (bool, error) {
				return false, fmt.Errorf("fail")
			}).
			And(Contains("qa"))("test", Args{})

		assert.NotNil(t, err)
		assert.False(t, result)
	})
}
