package test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestPostJSON(t *testing.T) {
	type jsonTestModel struct {
		Name string `json:"name"`
		OK   bool   `json:"ok"`
	}

	t.Run("should match specific json body fields", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(mocha.Postf("/test").
			Header("test", Equal("hello")).
			Body(
				JSONPath("name", Equal("dev")), JSONPath("ok", Equal(true))).
			Reply(mocha.OK()))

		req := testutil.PostJSON(m.URL()+"/test", &jsonTestModel{Name: "dev", OK: true})
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a struct", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		data := &jsonTestModel{OK: true, Name: "dev"}

		scoped := m.MustMock(mocha.Post(URLPath("/test")).
			Header("test", Equal("hello")).
			Body(EqualJSON(data)).
			Reply(mocha.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a map", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		data1 := map[string]interface{}{"ok": true, "name": "dev"}
		data2 := map[string]interface{}{"ok": true, "name": "dev"}

		scoped := m.MustMock(mocha.Post(URLPath("/test")).
			Header("test", Equal("hello")).
			Body(EqualJSON(data1)).
			Reply(mocha.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data2)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body when comparing a struct and a map", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		toMatch := map[string]interface{}{"name": "dev", "ok": true}
		data := jsonTestModel{Name: "dev", OK: true}

		scoped := m.MustMock(mocha.Post(URLPath("/test")).
			Header("test", Equal("hello")).
			Body(EqualJSON(toMatch)).
			Reply(mocha.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should not match when the given json is different than the incoming request body", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		body := map[string]interface{}{"ok": true, "name": "dev"}
		exp := map[string]interface{}{"ok": false, "name": "qa"}

		scoped := m.MustMock(mocha.Post(URLPath("/test")).
			Header("test", Equal("hello")).
			Body(EqualJSON(exp)).
			Reply(mocha.OK()))

		req := testutil.PostJSON(m.URL()+"/test", body)
		req.Header("test", "hello")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, mocha.StatusRequestDidNotMatch, res.StatusCode)
	})

	t.Run("should match null fields", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(mocha.Postf("/test").
			Body(
				JSONPath("name", Equal(nil)), JSONPath("ok", Equal(true))).
			Reply(mocha.OK()))

		req := testutil.Post(m.URL()+"/test", strings.NewReader(`{"name": null, "ok": true}`))
		req.Header("Content-Type", "application/json")

		res, err := req.Do()

		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestMalformedJSON_ShouldMatchOtherFieldsAndContinue(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(mocha.Postf("/test").
		Header("test", Equal("hello")).
		Reply(mocha.OK()))

	req := testutil.Post(m.URL()+"/test", strings.NewReader(`{"test": "malformed_json", "pass`))
	req.Header("test", "hello")
	req.Header("Content-Type", "application/json")

	res, err := req.Do()

	assert.NoError(t, err)
	assert.NoError(t, res.Body.Close())
	assert.True(t, scoped.HasBeenCalled())
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
