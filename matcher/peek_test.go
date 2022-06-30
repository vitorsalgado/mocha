package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeek(t *testing.T) {
	params := Args{}

	t.Run("should return peek error", func(t *testing.T) {
		actionErr := fmt.Errorf("action failed")
		result, err := Peek(EqualTo("test"), func(_ string) error { return actionErr })("test", params)

		assert.NotNil(t, err)
		assert.Equal(t, actionErr, err)
		assert.False(t, result)
	})

	t.Run("should execute action before returning provided matcher result", func(t *testing.T) {
		c := ""
		result, err := Peek(EqualTo("test"), func(v string) error { c = v; return nil })("test", params)

		assert.Nil(t, err)
		assert.True(t, result)
		assert.Equal(t, "test", c)
	})
}
