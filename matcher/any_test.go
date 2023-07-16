package matcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnyOf(t *testing.T) {
	t.Run("should return true if any of the given matchers returns true", func(t *testing.T) {
		result, err := Any(
			StrictEqual("test"),
			EqualIgnoreCase("dev"),
			ToLower(StrictEqual("TEST")),
			Contain("qa")).
			Match("test")
		require.Nil(t, err)
		require.True(t, result.Pass)
	})

	t.Run("should return false if all of the given matchers returns false", func(t *testing.T) {
		result, err := Any(
			StrictEqual("abc"),
			EqualIgnoreCase("def"),
			ToLower(StrictEqual("TEST")),
			Contain("dev")).
			Match("test")
		require.Nil(t, err)
		require.False(t, result.Pass)
	})

	t.Run("should return error and false when any error occurs", func(t *testing.T) {
		res, err := Any(
			StrictEqual("dev"),
			ToUpper(StrictEqual("none")),
			Func(func(v any) (bool, error) {
				return false, errors.New("boom")
			}),
			EqualIgnoreCase("qa")).
			Match("test")
		require.Error(t, err)
		require.False(t, res.Pass)
	})

	t.Run("no matchers", func(t *testing.T) {
		require.Panics(t, func() {
			Any()
		})
	})
}
