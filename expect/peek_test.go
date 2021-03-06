package expect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeek(t *testing.T) {
	params := Args{}

	t.Run("should return peek error", func(t *testing.T) {
		actionErr := fmt.Errorf("action failed")
		result, err := Peek(ToEqual("test"), func(_ any) error { return actionErr }).Matches("test", params)

		assert.NotNil(t, err)
		assert.Equal(t, actionErr, err)
		assert.False(t, result)
	})

	t.Run("should execute action before returning provided matcher result", func(t *testing.T) {
		c := ""
		result, err := Peek(ToEqual("test"), func(v any) error { c = v.(string); return nil }).Matches("test", params)

		assert.Nil(t, err)
		assert.True(t, result)
		assert.Equal(t, "test", c)
	})
}
