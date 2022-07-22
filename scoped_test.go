package mocha

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/core/_mocks"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
)

func TestScoped(t *testing.T) {
	m1 := core.NewMock()
	m2 := core.NewMock()
	m3 := core.NewMock()

	repo := core.NewStorage()
	repo.Save(m1)
	repo.Save(m2)
	repo.Save(m3)

	scoped := scope(repo, repo.FetchAll())

	assert.Equal(t, 3, len(scoped.ListAll()))
	assert.Equal(t, m1, scoped.Get(m1.ID))

	t.Run("should not return done when there is still pending mocks", func(t *testing.T) {
		fakeT := mocks.NewFakeNotifier()

		assert.False(t, scoped.Called())
		assert.Equal(t, 3, len(scoped.ListPending()))
		assert.True(t, scoped.IsPending())

		scoped.AssertCalled(fakeT)
		fakeT.AssertNumberOfCalls(t, "Errorf", 1)
	})

	t.Run("should return done when all mocks were called", func(t *testing.T) {
		fakeT := mocks.NewFakeNotifier()

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

		s1 := m.Mock(Get(expect.URLPath("/test1")).Reply(reply.OK()))
		s2 := m.Mock(
			Get(expect.URLPath("/test2")).Reply(reply.OK()),
			Get(expect.URLPath("/test3")).Reply(reply.OK()))

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
