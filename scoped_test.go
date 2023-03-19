package mocha

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestScoped(t *testing.T) {
	m1 := newMock()
	m2 := newMock()
	m3 := newMock()

	repo := newStore()
	repo.Save(m1)
	repo.Save(m2)
	repo.Save(m3)

	scoped := newScope(repo, []string{m1.ID, m2.ID, m3.ID})

	assert.Equal(t, 3, len(scoped.GetAll()))
	assert.Equal(t, m1, scoped.Get(m1.ID))
	assert.Nil(t, scoped.Get("unknown"))

	t.Run("should not return done when there is still pending mocks", func(t *testing.T) {
		ft := newFakeT()

		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, 3, len(scoped.GetPending()))
		assert.True(t, scoped.IsPending())

		scoped.AssertCalled(ft)
		ft.AssertNumberOfCalls(t, "Errorf", 1)
	})

	t.Run("should return done when all mocks were called", func(t *testing.T) {
		ft := newFakeT()

		m1.Inc()

		assert.False(t, scoped.HasBeenCalled())
		scoped.AssertCalled(ft)

		m2.Inc()
		m3.Inc()

		ft.AssertNumberOfCalls(t, "Errorf", 1)
		assert.True(t, scoped.AssertCalled(t))
		assert.True(t, scoped.HasBeenCalled())
		assert.True(t, scoped.AssertNumberOfCalls(t, 3))
		assert.False(t, scoped.AssertNumberOfCalls(ft, 1))
		assert.Equal(t, 0, len(scoped.GetPending()))
		assert.False(t, scoped.IsPending())
	})

	t.Run("should return total hits from store", func(t *testing.T) {
		assert.Equal(t, 3, scoped.Hits())
		scoped.AssertNumberOfCalls(t, 3)
	})

	t.Run("should clean all store associated with newScope when calling .Clean()", func(t *testing.T) {
		scoped.Clean()
		assert.Equal(t, 0, len(scoped.GetPending()))
		assert.True(t, scoped.AssertNumberOfCalls(t, 0))
		assert.False(t, scoped.IsPending())
	})

	t.Run("should only consider enabled store", func(t *testing.T) {
		m := New()
		m.MustStart()

		defer m.Close()

		s1 := m.MustMock(Get(matcher.URLPath("/test1")).Reply(OK()))
		s2 := m.MustMock(
			Get(matcher.URLPath("/test2")).Reply(OK()),
			Get(matcher.URLPath("/test3")).Reply(OK()))

		t.Run("initial state (enabled)", func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("%s/test1", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			res, err = http.Get(fmt.Sprintf("%s/test2", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			res, err = http.Get(fmt.Sprintf("%s/test3", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			assert.True(t, s1.HasBeenCalled())
			assert.True(t, s2.HasBeenCalled())
		})

		t.Run("disabled", func(t *testing.T) {
			s1.Disable()

			res, err := http.Get(fmt.Sprintf("%s/test1", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, StatusNoMatch, res.StatusCode)

			res, err = http.Get(fmt.Sprintf("%s/test2", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			res, err = http.Get(fmt.Sprintf("%s/test3", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 1, s1.Hits())

			s1.AssertNumberOfCalls(t, 1)
		})

		t.Run("enabling previously disabled", func(t *testing.T) {
			s1.Enable()

			res, err := http.Get(fmt.Sprintf("%s/test1", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 2, s1.Hits())

			s1.AssertNumberOfCalls(t, 2)
		})

		t.Run("disabling multiple", func(t *testing.T) {
			s2.Disable()

			res, err := http.Get(fmt.Sprintf("%s/test2", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, StatusNoMatch, res.StatusCode)

			res, err = http.Get(fmt.Sprintf("%s/test3", m.URL()))

			assert.NoError(t, err)
			assert.Equal(t, StatusNoMatch, res.StatusCode)
		})
	})
}

func TestScopedDelete(t *testing.T) {
	m1 := newMock()
	m2 := newMock()
	m3 := newMock()

	repo := newStore()
	repo.Save(m1)
	repo.Save(m2)
	repo.Save(m3)

	scoped := newScope(repo, []string{m1.ID, m2.ID, m3.ID})

	assert.True(t, scoped.Delete(m1.ID))
	assert.False(t, scoped.Delete("unknown"))

	assert.Nil(t, scoped.Get(m1.ID))
	assert.Nil(t, repo.Get(m1.ID))
}
