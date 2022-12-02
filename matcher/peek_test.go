package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeek(t *testing.T) {
	t.Run("should return peek error", func(t *testing.T) {
		actionErr := fmt.Errorf("action failed")
		result, err := Peek(Equal("test"), func(_ any) error { return actionErr }).Match("test")

		assert.NotNil(t, err)
		assert.Equal(t, actionErr, err)
		assert.False(t, result.OK)
	})

	t.Run("should execute action before returning provided matcher result", func(t *testing.T) {
		c := ""
		result, err := Peek(Equal("test"), func(v any) error { c = v.(string); return nil }).Match("test")

		assert.Nil(t, err)
		assert.True(t, result.OK)
		assert.Equal(t, "test", c)
	})
}