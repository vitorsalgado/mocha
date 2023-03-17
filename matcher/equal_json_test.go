package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToEqualJSON(t *testing.T) {
	t.Run("should return matcher error", func(t *testing.T) {
		c := make(chan bool, 1)
		body := map[string]interface{}{"ok": true, "name": "dev"}
		res, err := EqualJSON(c).Match(body)

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should match two equal json values", func(t *testing.T) {
		body := map[string]interface{}{"ok": true, "name": "dev"}
		res, err := Eqj(body).Match(body)

		require.NoError(t, err)
		require.True(t, res.Pass)
	})

	t.Run("should not match two different json values", func(t *testing.T) {
		a := map[string]interface{}{"ok": true, "name": "dev"}
		b := map[string]interface{}{"nok": true, "name": "dev"}
		res, err := Eqj(a).Match(b)

		require.NoError(t, err)
		require.False(t, res.Pass)
	})
}

func TestEqualJSONMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Eqj("").Name())
}
