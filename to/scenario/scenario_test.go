package scenario

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/to"
)

func TestScenario(t *testing.T) {
	t.Run("should init scenario as started", func(t *testing.T) {
		assert.True(t, newScenario("test").HasStarted())
	})

	t.Run("should only create scenario if needed", func(t *testing.T) {
		store := NewStore()
		store.CreateNewIfNeeded("scenario-1")

		s, ok := store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.True(t, s.HasStarted())

		s.State = "another-state"
		store.Save(s)

		store.CreateNewIfNeeded("scenario-1")

		s, ok = store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.False(t, s.HasStarted())
		assert.Equal(t, s.State, "another-state")
	})
}

func TestScenarioConditions(t *testing.T) {
	store := NewStore()
	p := params.New()
	p.Set(BuiltInParamStore, store)
	args := to.Args{Params: p}

	t.Run("should return true when scenario is not started and also not found", func(t *testing.T) {
		m := Scenario[any]("test", "required", "newScenario")
		res, err := m.Matches(nil, args)

		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return false when scenario exists but it is not in the required state", func(t *testing.T) {
		store.CreateNewIfNeeded("hi")

		m := Scenario[any]("hi", "required", "newScenario")
		res, err := m.Matches(nil, args)

		assert.Nil(t, err)
		assert.False(t, res)
	})
}
