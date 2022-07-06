package matchers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBothMatcher(t *testing.T) {
	t.Parallel()

	t.Run("should return false when only left matcher evaluates to true", func(t *testing.T) {
		result, err := Both(EqualTo("test")).And(Contains("qa")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when only right matcher evaluates to true", func(t *testing.T) {
		result, err := Both(EqualTo("qa")).And(Contains("tes")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when both matchers evaluates to true", func(t *testing.T) {
		result, err := Both(EqualTo("test")).And(Contains("te")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when only left matcher evaluates to true", func(t *testing.T) {
		result, err := Both(EqualTo("test")).And(Contains("qa")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when matchers throws errors", func(t *testing.T) {
		result, err := Both(
			Fn(func(_ string, _ Args) (bool, error) {
				return false, fmt.Errorf("fail")
			})).
			And(Contains("qa")).Matches("test", emptyArgs())

		assert.NotNil(t, err)
		assert.False(t, result)
	})
}
