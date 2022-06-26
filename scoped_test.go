package mocha

import (
	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/mock"
	"testing"
)

func TestScoped(t *testing.T) {
	repo := mock.NewStorage()
	repo.Save(mock.New())
	repo.Save(mock.New())
	repo.Save(mock.New())

	scoped := Scope(repo, repo.FetchAll())

	t.Run("should not return done when there is still pending mocks", func(t *testing.T) {
		assert.False(t, scoped.IsDone())
		assert.Equal(t, 3, len(scoped.Pending()))
	})

	t.Run("should return done when all mocks were called", func(t *testing.T) {
		m := repo.FetchByID(1)
		m.Hit()

		assert.False(t, scoped.IsDone())
		assert.NotNil(t, scoped.Done())

		m = repo.FetchByID(2)
		m.Hit()

		m = repo.FetchByID(3)
		m.Hit()

		assert.True(t, scoped.IsDone())
		assert.Nil(t, scoped.Done())
		assert.Equal(t, 0, len(scoped.Pending()))
	})

	t.Run("should clean all mocks associated with scope when calling .Clean()", func(t *testing.T) {
		scoped.Clean()
		assert.Equal(t, 0, len(scoped.Pending()))
		assert.False(t, scoped.IsPending())
	})
}
