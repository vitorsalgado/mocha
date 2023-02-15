package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeek(t *testing.T) {
	t.Run("should return peek error", func(t *testing.T) {
		actionErr := fmt.Errorf("action failed")
		result, err := Peek(StrictEqual("test"), func(_ any) error { return actionErr }).Match("test")

		require.Error(t, err)
		assert.Equal(t, actionErr, err)
		assert.Nil(t, result)
	})

	t.Run("should execute action before returning provided matcher result", func(t *testing.T) {
		c := ""
		result, err := Peek(StrictEqual("test"), func(v any) error { c = v.(string); return nil }).Match("test")

		require.NoError(t, err)
		assert.True(t, result.Pass)
		assert.Equal(t, "test", c)
	})
}
