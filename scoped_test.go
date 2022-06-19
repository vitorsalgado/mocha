package mocha

import (
	"testing"

	"github.com/vitorsalgado/mocha/internal/assert"
)

func TestScoped(t *testing.T) {
	repo := NewMockStore()
	repo.Save(&Mock{ID: 1})
	repo.Save(&Mock{ID: 2})
	repo.Save(&Mock{ID: 3})

	ids := []int32{1, 2}
	scoped := NewScoped(repo, ids)

	t.Run("should not return done when there is still pending mocks", func(t *testing.T) {
		assert.False(t, scoped.IsDone())
		assert.Equal(t, 2, len(scoped.Pending()))
	})

	t.Run("should return done when all mocks were called", func(t *testing.T) {
		m := repo.GetByID(1)
		m.Hit()

		repo.Save(&m)

		assert.False(t, scoped.IsDone())
		assert.NotNil(t, scoped.Done())

		m = repo.GetByID(2)
		m.Hit()

		repo.Save(&m)

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
