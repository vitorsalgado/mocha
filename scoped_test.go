package mocha

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/internal/testmocks"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestScoped(t *testing.T) {
	m1 := newMock()
	m2 := newMock()
	m3 := newMock()

	repo := newStorage()
	repo.Save(m1)
	repo.Save(m2)
	repo.Save(m3)

	scoped := scope(repo, repo.FetchAll())

	assert.Equal(t, 3, len(scoped.ListAll()))
	assert.Equal(t, m1, scoped.Get(m1.ID))

	t.Run("should not return done when there is still pending mocks", func(t *testing.T) {
		fakeT := testmocks.NewFakeNotifier()

		assert.False(t, scoped.Called())
		assert.Equal(t, 3, len(scoped.ListPending()))
		assert.True(t, scoped.IsPending())

		scoped.AssertCalled(fakeT)
		fakeT.AssertNumberOfCalls(t, "Errorf", 1)
	})

	t.Run("should return done when all mocks were called", func(t *testing.T) {
		fakeT := testmocks.NewFakeNotifier()

		m1.Hit()

		assert.False(t, scoped.Called())
		scoped.AssertCalled(fakeT)

		m2.Hit()
		m3.Hit()

		fakeT.AssertNumberOfCalls(t, "Errorf", 1)
		assert.True(t, scoped.AssertCalled(t))
		assert.True(t, scoped.Called())
		assert.Equal(t, 0, len(scoped.ListPending()))
		assert.False(t, scoped.IsPending())
	})

	t.Run("should return total hits from mocks", func(t *testing.T) {
		assert.Equal(t, 3, scoped.Hits())
	})

	t.Run("should clean all mocks associated with scope when calling .Clean()", func(t *testing.T) {
		scoped.Clean()
		assert.Equal(t, 0, len(scoped.ListPending()))
		assert.False(t, scoped.IsPending())
	})

	t.Run("should only consider enabled mocks", func(t *testing.T) {
		m := New(t)
		m.Start()

		defer m.Close()

		s1 := m.AddMocks(Get(matcher.URLPath("/test1")).Reply(reply.OK()))
		s2 := m.AddMocks(
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

			assert.True(t, s1.Called())
			assert.True(t, s2.Called())
		})

		t.Run("disabled", func(t *testing.T) {
			s1.Disable()

			req := testutil.Get(fmt.Sprintf("%s/test1", m.URL()))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusTeapot, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test2", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 1, s1.Hits())
		})

		t.Run("enabling previously disabled", func(t *testing.T) {
			s1.Enable()

			req := testutil.Get(fmt.Sprintf("%s/test1", m.URL()))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 2, s1.Hits())
		})

		t.Run("disabling multiple", func(t *testing.T) {
			s2.Disable()

			req := testutil.Get(fmt.Sprintf("%s/test2", m.URL()))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusTeapot, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.URL()))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusTeapot, res.StatusCode)
		})
	})
}
