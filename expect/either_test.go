package expect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEither(t *testing.T) {
	t.Parallel()

	t.Run("should return true when only left matcher evaluates to true", func(t *testing.T) {
		result, err := Either(ToEqual("test")).Or(ToContain("qa")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when only right matcher evaluates to true", func(t *testing.T) {
		result, err := Either(ToEqual("qa")).Or(ToContain("tes")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when both matchers evaluates to true", func(t *testing.T) {
		result, err := Either(ToEqual("test")).Or(ToContain("te")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when both evaluates to false", func(t *testing.T) {
		result, err := Either(ToEqual("dev")).Or(ToContain("qa")).Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when matchers throws errors", func(t *testing.T) {
		result, err := Either(
			Func(func(_ any, _ Args) (bool, error) {
				return false, fmt.Errorf("fail")
			})).
			Or(ToContain("qa")).
			Matches("test", emptyArgs())

		assert.NotNil(t, err)
		assert.False(t, result)
	})
}
