package matcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFn(t *testing.T) {
	m := Func(func(v any) (bool, error) {
		return v.(string) == "ok", nil
	})

	t.Run("pass", func(t *testing.T) {
		res, err := m.Match("ok")
		require.NoError(t, err)
		require.True(t, res.Pass)
	})

	t.Run("no pass", func(t *testing.T) {
		res, err := m.Match("nok")
		require.NoError(t, err)
		require.False(t, res.Pass)
	})
}

func TestFnError(t *testing.T) {
	res, err := Func(
		func(_ any) (bool, error) {
			return false, errors.New("boom")
		}).
		Match("ok")

	require.Error(t, err)
	require.Nil(t, res)
}

func TestFuncMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Func(nil).Name())
}
