package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Run("should return true when value is contained in the matcher argument", func(t *testing.T) {
		result, err := ToContain("world").Matches("hello world", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})
	t.Run("should return false when value is not contained in the matcher argument", func(t *testing.T) {
		result, err := ToContain("dev").Matches("hello world", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when expected value is contained in the given slice", func(t *testing.T) {
		result, err := ToContain("dev").Matches([]string{"dev", "qa"}, emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})
	t.Run("should return false when expected value is not contained in the given slice", func(t *testing.T) {
		result, err := ToContain("po").Matches([]string{"dev", "qa"}, emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when expected value is a key present in the given map", func(t *testing.T) {
		result, err := ToContain("dev").Matches(map[string]string{"dev": "ok", "qa": "nok"}, emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})
	t.Run("should return false when expected value is a key not present in the given map", func(t *testing.T) {
		result, err := ToContain("unknown").Matches(map[string]string{"dev": "ok", "qa": "nok"}, emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})
}
