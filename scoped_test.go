package mocha

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/reply"
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

	t.Run("should only consider enalbed mocks", func(t *testing.T) {
		m := ForTest(t)
		m.Start()

		s1 := m.Mock(Get(matcher.URLPath("/test1")).Reply(reply.OK()))
		s2 := m.Mock(
			Get(matcher.URLPath("/test2")).Reply(reply.OK()),
			Get(matcher.URLPath("/test3")).Reply(reply.OK()))

		t.Run("initial state (enabled)", func(t *testing.T) {
			req := testutil.Get(fmt.Sprintf("%s/test1", m.Server.URL))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test2", m.Server.URL))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.Server.URL))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			assert.True(t, s1.IsDone())
			assert.True(t, s2.IsDone())
		})

		t.Run("disabled", func(t *testing.T) {
			s1.Disable()

			req := testutil.Get(fmt.Sprintf("%s/test1", m.Server.URL))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusTeapot, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test2", m.Server.URL))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.Server.URL))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 1, s1.Hits())
		})

		t.Run("enabling previously disabled", func(t *testing.T) {
			s1.Enable()

			req := testutil.Get(fmt.Sprintf("%s/test1", m.Server.URL))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 2, s1.Hits())
		})

		t.Run("disabling multiple", func(t *testing.T) {
			s2.Disable()

			req := testutil.Get(fmt.Sprintf("%s/test2", m.Server.URL))
			res, err := req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusTeapot, res.StatusCode)

			req = testutil.Get(fmt.Sprintf("%s/test3", m.Server.URL))
			res, err = req.Do()

			assert.NoError(t, err)
			assert.Equal(t, http.StatusTeapot, res.StatusCode)
		})
	})
}
