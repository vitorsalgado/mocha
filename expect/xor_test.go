package expect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXor(t *testing.T) {
	m := XOR(ToContain("dev"), ToContain("test"))

	t.Run("should return true when left condition matches", func(t *testing.T) {
		res, err := m.Match("dev")
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return true right condition matches", func(t *testing.T) {
		res, err := m.Match("test")
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return false both condition matches", func(t *testing.T) {
		res, err := m.Match("dev-test")
		assert.Nil(t, err)
		assert.False(t, res)
	})
}

func TestXorError(t *testing.T) {
	t.Run("should return error from right matcher and return false", func(t *testing.T) {
		m := XOR(ToContain("dev"), Func(func(_ any) (bool, error) {
			return false, fmt.Errorf("fail")
		}))

		res, err := m.Match("dev")

		assert.Error(t, err)
		assert.False(t, res)
	})

	t.Run("should return error from left matcher and return false", func(t *testing.T) {
		m := XOR(
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			ToContain("dev"))

		res, err := m.Match("dev")

		assert.Error(t, err)
		assert.False(t, res)
	})

	t.Run("should return error when both matchers fails", func(t *testing.T) {
		m := XOR(
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail First")
			}),
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail second")
			}))

		res, err := m.Match("nothing")

		assert.Error(t, err)
		assert.False(t, res)
	})
}
