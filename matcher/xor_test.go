package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXor(t *testing.T) {
	m := XOR(Contain("dev"), Contain("test"))

	t.Run("should return true when left condition matches", func(t *testing.T) {
		res, err := m.Match("dev")
		assert.Nil(t, err)
		assert.True(t, res.Pass)
	})

	t.Run("should return true right condition matches", func(t *testing.T) {
		res, err := m.Match("test")
		assert.Nil(t, err)
		assert.True(t, res.Pass)
	})

	t.Run("should return false both condition matches", func(t *testing.T) {
		res, err := m.Match("dev-test")
		assert.Nil(t, err)
		assert.False(t, res.Pass)
	})
}

func TestXorError(t *testing.T) {
	t.Run("should return error from right matcher and return false", func(t *testing.T) {
		m := XOR(Contain("dev"), Func(func(_ any) (bool, error) {
			return false, fmt.Errorf("fail")
		}))

		res, err := m.Match("dev")

		assert.Error(t, err)
		assert.False(t, res.Pass)
	})

	t.Run("should return error from left matcher and return false", func(t *testing.T) {
		m := XOR(
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			Contain("dev"))

		res, err := m.Match("dev")

		assert.Error(t, err)
		assert.False(t, res.Pass)
	})

	t.Run("should return error when both matchers fails", func(t *testing.T) {
		m := XOR(
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail first")
			}),
			Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail second")
			}))

		res, err := m.Match("nothing")

		assert.Error(t, err)
		assert.False(t, res.Pass)
	})
}
