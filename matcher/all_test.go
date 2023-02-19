package matcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllOf(t *testing.T) {
	t.Run("should return true when all matchers evaluates to true", func(t *testing.T) {
		result, err := All(
			StrictEqual("test"),
			EqualIgnoreCase("TEST"),
			ToUpper(StrictEqual("TEST")),
			Contain("tes")).
			Match("test")
		require.NoError(t, err)
		require.True(t, result.Pass)
	})

	t.Run("should return false when just one matcher evaluates to false", func(t *testing.T) {
		result, err := All(
			StrictEqual("test"),
			EqualIgnoreCase("dev"),
			ToUpper(StrictEqual("TEST")),
			Contain("tes")).
			Match("test")
		require.NoError(t, err)
		require.False(t, result.Pass)
	})

	t.Run("should return false when all matchers evaluates to false", func(t *testing.T) {
		result, err := All(
			StrictEqual("dev"),
			EqualIgnoreCase("qa"),
			ToUpper(StrictEqual("none")),
			Contain("blah")).
			Match("test")
		require.NoError(t, err)
		require.False(t, result.Pass)
	})

	t.Run("should return false when an error occurs", func(t *testing.T) {
		_, err := All(
			StrictEqual("dev"),
			EqualIgnoreCase("qa"),
			ToUpper(StrictEqual("none")),
			Func(func(v any) (bool, error) {
				return false, errors.New("boom")
			})).
			Match("test")
		require.Error(t, err)
	})

	t.Run("no matchers", func(t *testing.T) {
		require.Panics(t, func() {
			All()
		})
	})
}
