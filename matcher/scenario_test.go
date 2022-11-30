package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScenario(t *testing.T) {
	t.Run("should init scenario as started", func(t *testing.T) {
		assert.True(t, newScenarioState("test").hasStarted())
	})

	t.Run("should only create scenario if needed", func(t *testing.T) {
		store := NewScenarioStorage()
		store.CreateNewIfNeeded("scenario-1")

		s, ok := store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.True(t, s.hasStarted())

		s.State = "another-state"
		store.Save(s)

		store.CreateNewIfNeeded("scenario-1")

		s, ok = store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.False(t, s.hasStarted())
		assert.Equal(t, s.State, "another-state")
	})
}
