package matchers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXor(t *testing.T) {
	m := XOR(Contains("dev"), Contains("test"))

	t.Run("should return true when left condition matches", func(t *testing.T) {
		res, err := m.Matches("dev", Args{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return true right condition matches", func(t *testing.T) {
		res, err := m.Matches("test", Args{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return false both condition matches", func(t *testing.T) {
		res, err := m.Matches("dev-test", Args{})
		assert.Nil(t, err)
		assert.False(t, res)
	})
}

func TestXorError(t *testing.T) {
	t.Run("should return error from right matcher and return false", func(t *testing.T) {
		m := XOR(Contains("dev"), Fn(func(_ string, _ Args) (bool, error) {
			return false, fmt.Errorf("fail")
		}))

		res, err := m.Matches("dev", Args{})

		assert.Error(t, err)
		assert.False(t, res)
	})

	t.Run("should return error from left matcher and return false", func(t *testing.T) {
		m := XOR(
			Fn(func(_ string, _ Args) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			Contains("dev"))

		res, err := m.Matches("dev", Args{})

		assert.Error(t, err)
		assert.False(t, res)
	})

	t.Run("should return error when both matchers fails", func(t *testing.T) {
		m := XOR(
			Fn(func(_ string, _ Args) (bool, error) {
				return false, fmt.Errorf("fail firts")
			}),
			Fn(func(_ string, _ Args) (bool, error) {
				return false, fmt.Errorf("fail second")
			}))

		res, err := m.Matches("nothing", Args{})

		assert.Error(t, err)
		assert.False(t, res)
	})
}
