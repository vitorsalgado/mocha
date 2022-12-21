package mocha

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestScoped(t *testing.T) {
	m1 := newMock()
	m2 := newMock()
	m3 := newMock()

	repo := newStore()
	repo.Save(m1)
	repo.Save(m2)
	repo.Save(m3)

	scoped := scope(repo, repo.GetAll())

	assert.Equal(t, 3, len(scoped.GetAll()))
	assert.Equal(t, m1, scoped.Get(m1.ID))

	t.Run("should not return done when there is still pending store", func(t *testing.T) {
		fakeT := NewFakeNotifier()

		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, 3, len(scoped.GetPending()))
		assert.True(t, scoped.IsPending())

		scoped.AssertCalled(fakeT)
		fakeT.AssertNumberOfCalls(t, "Errorf", 1)
	})

	t.Run("should return done when all store were called", func(t *testing.T) {
		fakeT := NewFakeNotifier()

		m1.Inc()

		assert.False(t, scoped.HasBeenCalled())
		scoped.AssertCalled(fakeT)

		m2.Inc()
		m3.Inc()

		fakeT.AssertNumberOfCalls(t, "Errorf", 1)
		assert.True(t, scoped.AssertCalled(t))
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, 0, len(scoped.GetPending()))
		assert.False(t, scoped.IsPending())
	})

	t.Run("should return total hits from store", func(t *testing.T) {
		assert.Equal(t, 3, scoped.Hits())
		scoped.AssertNumberOfCalls(t, 3)
	})

	t.Run("should clean all store associated with scope when calling .Clean()", func(t *testing.T) {
		scoped.Clean()
		assert.Equal(t, 0, len(scoped.GetPending()))
		assert.False(t, scoped.IsPending())
	})

	t.Run("should only consider enabled store", func(t *testing.T) {
		m := NewWithT(t)
		m.MustStart()

		defer m.Close()

		s1 := m.MustMock(Get(matcher.URLPath("/test1")).Reply(reply.OK()))
		s2 := m.MustMock(
			Get(matcher.URLPath("/test2")).Reply(reply.OK()),
			Get(matcher.URLPath("/test3")).Reply(reply.OK()))

		t.Run("initial state (enabled)", func(t *testing.T) {
			req := testutil.Get(fmt.Sprintf("%s/test1", m.URL()))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test2", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			assert.True(t, s1.HasBeenCalled())
			assert.True(t, s2.HasBeenCalled())
		})

		t.Run("disabled", func(t *testing.T) {
			s1.Disable()

			req := testutil.Get(fmt.Sprintf("%s/test1", m.URL()))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, StatusRequestDidNotMatch, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test2", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 1, s1.Hits())

			s1.AssertNumberOfCalls(t, 1)
		})

		t.Run("enabling previously disabled", func(t *testing.T) {
			s1.Enable()

			req := testutil.Get(fmt.Sprintf("%s/test1", m.URL()))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 2, s1.Hits())

			s1.AssertNumberOfCalls(t, 2)
		})

		t.Run("disabling multiple", func(t *testing.T) {
			s2.Disable()

			req := testutil.Get(fmt.Sprintf("%s/test2", m.URL()))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, StatusRequestDidNotMatch, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, StatusRequestDidNotMatch, res.StatusCode)
		})
	})
}
