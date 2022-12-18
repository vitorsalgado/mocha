package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToEqualJSON(t *testing.T) {
	t.Run("should return matcher error", func(t *testing.T) {
		c := make(chan bool, 1)
		body := map[string]interface{}{"ok": true, "name": "dev"}
		res, err := EqualJSON(c).Match(body)

		assert.Error(t, err)
		assert.False(t, res.Pass)
	})
}
