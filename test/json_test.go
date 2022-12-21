package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type jsonTestModel struct {
	Name string `json:"name"`
	OK   bool   `json:"ok"`
}

func TestPostJSON(t *testing.T) {
	t.Run("should match specific json body fields", func(t *testing.T) {
		m := mocha.NewWithT(t)
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(mocha.Postf("/test").
			Header("test", matcher.Equal("hello")).
			Body(
				matcher.JSONPath("name", matcher.Equal("dev")), matcher.JSONPath("ok", matcher.Equal(true))).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", &jsonTestModel{Name: "dev", OK: true})
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a struct", func(t *testing.T) {
		m := mocha.NewWithT(t)
		m.MustStart()

		data := &jsonTestModel{OK: true, Name: "dev"}

		scoped := m.MustMock(mocha.Post(matcher.URLPath("/test")).
			Header("test", matcher.Equal("hello")).
			Body(matcher.EqualJSON(data)).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a map", func(t *testing.T) {
		m := mocha.NewWithT(t)
		m.MustStart()

		data1 := map[string]interface{}{"ok": true, "name": "dev"}
		data2 := map[string]interface{}{"ok": true, "name": "dev"}

		scoped := m.MustMock(mocha.Post(matcher.URLPath("/test")).
			Header("test", matcher.Equal("hello")).
			Body(matcher.EqualJSON(data1)).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data2)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body when comparing a struct and a map", func(t *testing.T) {
		m := mocha.NewWithT(t)
		m.MustStart()

		toMatch := map[string]interface{}{"name": "dev", "ok": true}
		data := jsonTestModel{Name: "dev", OK: true}

		scoped := m.MustMock(mocha.Post(matcher.URLPath("/test")).
			Header("test", matcher.Equal("hello")).
			Body(matcher.EqualJSON(toMatch)).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should not match when the given json is different than the incoming request body", func(t *testing.T) {
		m := mocha.NewWithT(t)
		m.MustStart()

		body := map[string]interface{}{"ok": true, "name": "dev"}
		exp := map[string]interface{}{"ok": false, "name": "qa"}

		scoped := m.MustMock(mocha.Post(matcher.URLPath("/test")).
			Header("test", matcher.Equal("hello")).
			Body(matcher.EqualJSON(exp)).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", body)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, mocha.StatusRequestDidNotMatch, res.StatusCode)
	})
}
