package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBothMatcher(t *testing.T) {
	t.Run("should return false when only left matcher evaluates to true", func(t *testing.T) {
		result, err := Both(Equal("test")).And(Contain("qa")).Match("test")
		assert.Nil(t, err)
		assert.False(t, result.OK)
	})

	t.Run("should return false when only right matcher evaluates to true", func(t *testing.T) {
		result, err := Both(Equal("qa")).And(Contain("tes")).Match("test")
		assert.Nil(t, err)
		assert.False(t, result.OK)
	})

	t.Run("should return true when both matchers evaluates to true", func(t *testing.T) {
		result, err := Both(Equal("test")).And(Contain("te")).Match("test")
		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return false when only left matcher evaluates to true", func(t *testing.T) {
		result, err := Both(Equal("test")).And(Contain("qa")).Match("test")
		assert.Nil(t, err)
		assert.False(t, result.OK)
	})

	t.Run("should return false when matchers throws errors", func(t *testing.T) {
		result, err := Both(
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			})).
			And(Contain("qa")).Match("test")

		assert.NotNil(t, err)
		assert.False(t, result.OK)
	})
}