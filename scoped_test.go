package mocha

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/mock"
)

func TestScoped(t *testing.T) {
	m1 := mock.New()
	m2 := mock.New()
	m3 := mock.New()

	repo := mock.NewStorage()
	repo.Save(m1)
	repo.Save(m2)
	repo.Save(m3)

	scoped := Scope(repo, repo.FetchAll())

	t.Run("should not return done when there is still pending mocks", func(t *testing.T) {
		assert.False(t, scoped.IsDone())
		assert.Equal(t, 3, len(scoped.Pending()))
	})

	t.Run("should return done when all mocks were called", func(t *testing.T) {
		m1.Hit()

		assert.False(t, scoped.IsDone())
		assert.NotNil(t, scoped.Done())

		m2.Hit()
		m3.Hit()

		assert.True(t, scoped.IsDone())
		assert.Nil(t, scoped.Done())
		assert.Equal(t, 0, len(scoped.Pending()))
	})

	t.Run("should return total hits from mocks", func(t *testing.T) {
		assert.Equal(t, 3, scoped.Hits())
	})

	t.Run("should clean all mocks associated with scope when calling .Clean()", func(t *testing.T) {
		scoped.Clean()
		assert.Equal(t, 0, len(scoped.Pending()))
		assert.False(t, scoped.IsPending())
	})
}
