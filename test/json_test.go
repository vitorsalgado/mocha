package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
)

type jsonTestModel struct {
	Name string `json:"name"`
	OK   bool   `json:"ok"`
}

func TestPostJSON(t *testing.T) {
	t.Run("should match specific json body fields", func(t *testing.T) {
		m := mocha.New(t)
		m.Start()

		scoped := m.Mock(mocha.Post(expect.URLPath("/test")).
			Header("test", expect.ToEqual("hello")).
			Body(
				expect.JSONPath("name", expect.ToEqual("dev")), expect.JSONPath("ok", expect.ToEqual(true))).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", &jsonTestModel{Name: "dev", OK: true})
		req.Header("test", "hello")

		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		defer res.Body.Close()

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a struct", func(t *testing.T) {
		m := mocha.New(t)
		m.Start()

		data := &jsonTestModel{OK: true, Name: "dev"}

		scoped := m.Mock(mocha.Post(expect.URLPath("/test")).
			Header("test", expect.ToEqual("hello")).
			Body(expect.ToEqualJSON(data)).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		defer res.Body.Close()

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a map", func(t *testing.T) {
		m := mocha.New(t)
		m.Start()

		data := map[string]interface{}{
			"ok":   true,
			"name": "dev",
		}

		scoped := m.Mock(mocha.Post(expect.URLPath("/test")).
			Header("test", expect.ToEqual("hello")).
			Body(expect.ToEqualJSON(data)).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		defer res.Body.Close()

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body when comparing a struct and a map", func(t *testing.T) {
		m := mocha.New(t)
		m.Start()

		toMatch := map[string]interface{}{
			"name": "dev",
			"ok":   true,
		}

		data := jsonTestModel{Name: "dev", OK: true}

		scoped := m.Mock(mocha.Post(expect.URLPath("/test")).
			Header("test", expect.ToEqual("hello")).
			Body(expect.ToEqualJSON(toMatch)).
			Reply(reply.OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		defer res.Body.Close()

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
