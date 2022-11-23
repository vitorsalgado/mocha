package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Run("should return true when value is contained in the matcher argument", func(t *testing.T) {
		result, err := ToContain("world").Match("hello world")
		assert.Nil(t, err)
		assert.True(t, result)
	})
	t.Run("should return false when value is not contained in the matcher argument", func(t *testing.T) {
		result, err := ToContain("dev").Match("hello world")
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when expected value is contained in the given slice", func(t *testing.T) {
		result, err := ToContain("dev").Match([]string{"dev", "qa"})
		assert.Nil(t, err)
		assert.True(t, result)
	})
	t.Run("should return false when expected value is not contained in the given slice", func(t *testing.T) {
		result, err := ToContain("po").Match([]string{"dev", "qa"})
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when expected value is a key present in the given map", func(t *testing.T) {
		result, err := ToContain("dev").Match(map[string]string{"dev": "ok", "qa": "nok"})
		assert.Nil(t, err)
		assert.True(t, result)
	})
	t.Run("should return false when expected value is a key not present in the given map", func(t *testing.T) {
		result, err := ToContain("unknown").Match(map[string]string{"dev": "ok", "qa": "nok"})
		assert.Nil(t, err)
		assert.False(t, result)
	})
}
