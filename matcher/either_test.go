package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEither(t *testing.T) {
	t.Run("should return true when only left matcher evaluates to true", func(t *testing.T) {
		result, err := Either(Equal("test"), Contain("qa")).Match("test")
		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return true when only right matcher evaluates to true", func(t *testing.T) {
		result, err := Either(Equal("qa"), Contain("tes")).Match("test")
		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return true when both matchers evaluates to true", func(t *testing.T) {
		result, err := Either(Equal("test"), Contain("te")).Match("test")
		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return false when both evaluates to false", func(t *testing.T) {
		result, err := Either(Equal("dev"), Contain("qa")).Match("test")
		assert.Nil(t, err)
		assert.False(t, result.OK)
	})

	t.Run("should return false when matchers throws errors", func(t *testing.T) {
		result, err := Either(
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			Contain("qa")).
			Match("test")

		assert.NotNil(t, err)
		assert.False(t, result.OK)
	})
}
