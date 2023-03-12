package mfeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScenario(t *testing.T) {
	t.Run("should init scenario as started", func(t *testing.T) {
		assert.True(t, newScenarioState("test").hasStarted())
	})

	t.Run("should only create scenario if needed", func(t *testing.T) {
		store := NewScenarioStore()
		store.createNewIfNeeded("scenario-1")

		s, ok := store.fetchByName("scenario-1")
		assert.True(t, ok)
		assert.True(t, s.hasStarted())

		s.state = "another-state"
		store.createNewIfNeeded("scenario-1")

		s, ok = store.fetchByName("scenario-1")
		assert.True(t, ok)
		assert.False(t, s.hasStarted())
		assert.Equal(t, s.state, "another-state")
	})
}
