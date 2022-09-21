package test

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestHTTPMethods(t *testing.T) {
	m := mocha.New(t)
	m.Start()

	t.Run("should mock GET", func(t *testing.T) {
		scoped := m.AddMocks(
			mocha.Get(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 1, scoped.Hits())

		other, err := testutil.Post(m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
	})

	t.Run("should mock POST", func(t *testing.T) {
		scoped := m.AddMocks(
			mocha.Post(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.Post(m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
	})

	t.Run("should mock PUT", func(t *testing.T) {
		scoped := m.AddMocks(
			mocha.Put(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.NewRequest(http.MethodPut, m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
	})

	t.Run("should mock DELETE", func(t *testing.T) {
		scoped := m.AddMocks(
			mocha.Delete(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.NewRequest(http.MethodDelete, m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
	})

	t.Run("should mock PATCH", func(t *testing.T) {
		scoped := m.AddMocks(
			mocha.Patch(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.NewRequest(http.MethodPatch, m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
	})

	t.Run("should mock HEAD", func(t *testing.T) {
		scoped := m.AddMocks(
			mocha.Head(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.NewRequest(http.MethodHead, m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
	})

	t.Run("should mock OPTIONS", func(t *testing.T) {
		scoped := m.AddMocks(
			mocha.Options(expect.URLPath("/test")).
				Reply(reply.OK()))

		defer scoped.Clean()

		res, err := testutil.Get(m.URL() + "/test").Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.False(t, scoped.Called())

		other, err := testutil.NewRequest(http.MethodOptions, m.URL()+"/test", nil).Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, other.StatusCode)
		assert.Equal(t, 1, scoped.Hits())
		assert.True(t, scoped.Called())
	})
}
