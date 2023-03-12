package test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestPostJSON(t *testing.T) {
	type level2 struct {
		Active *bool   `json:"active"`
		Value  float64 `json:"value"`
	}

	type level1 struct {
		Year      int       `json:"year"`
		Timestamp time.Time `json:"timestamp"`
		Day       int8      `json:"day"`
		Level2    level2    `json:"level2"`
	}

	type jsonTestModel struct {
		Name      string   `json:"name"`
		OK        bool     `json:"ok"`
		Unordered []string `json:"unordered"`
		Ordered   []string `json:"ordered"`
		Diff      []any    `json:"diff"`
		Level1    level1   `json:"level1"`
	}

	t.Run("should match specific json body fields", func(t *testing.T) {
		m := New()
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(Postf("/test").
			Header("test", StrictEqual("hello")).
			Body(
				JSONPath("name", StrictEqual("dev")), JSONPath("ok", StrictEqual(true))).
			Reply(OK()))

		req := testutil.PostJSON(m.URL()+"/test", &jsonTestModel{Name: "dev", OK: true})
		req.Header("test", "hello")

		res, err := req.Do()

		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a struct", func(t *testing.T) {
		m := New()
		m.MustStart()

		defer m.Close()

		data := &jsonTestModel{OK: true, Name: "dev"}

		scoped := m.MustMock(Post(URLPath("/test")).
			Header("test", StrictEqual("hello")).
			Body(EqualJSON(data)).
			Reply(OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()

		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body using a map", func(t *testing.T) {
		m := New()
		m.MustStart()

		defer m.Close()

		data1 := map[string]interface{}{"ok": true, "name": "dev"}
		data2 := map[string]interface{}{"ok": true, "name": "dev"}

		scoped := m.MustMock(Post(URLPath("/test")).
			Header("test", StrictEqual("hello")).
			Body(EqualJSON(data1)).
			Reply(OK()))

		req := testutil.PostJSON(m.URL()+"/test", data2)
		req.Header("test", "hello")

		res, err := req.Do()

		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should match entire body when comparing a struct and a map", func(t *testing.T) {
		type model struct {
			Name string `json:"name"`
			OK   bool   `json:"ok"`
		}

		m := New()
		m.MustStart()

		defer m.Close()

		toMatch := map[string]interface{}{"name": "dev", "ok": true}
		data := model{Name: "dev", OK: true}

		scoped := m.MustMock(Post(URLPath("/test")).
			Header("test", StrictEqual("hello")).
			Body(EqualJSON(toMatch)).
			Reply(OK()))

		req := testutil.PostJSON(m.URL()+"/test", data)
		req.Header("test", "hello")

		res, err := req.Do()

		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should not match when the given json is different than the incoming request body", func(t *testing.T) {
		m := New()
		m.MustStart()

		defer m.Close()

		body := map[string]interface{}{"ok": true, "name": "dev"}
		exp := map[string]interface{}{"ok": false, "name": "qa"}

		scoped := m.MustMock(Post(URLPath("/test")).
			Header("test", StrictEqual("hello")).
			Body(EqualJSON(exp)).
			Reply(OK()))

		req := testutil.PostJSON(m.URL()+"/test", body)
		req.Header("test", "hello")

		res, err := req.Do()

		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, StatusNoMatch, res.StatusCode)
	})

	t.Run("should match null fields", func(t *testing.T) {
		m := New()
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(Postf("/test").
			Body(
				JSONPath("name", StrictEqual(nil)), JSONPath("ok", StrictEqual(true))).
			Reply(OK()))

		req := testutil.Post(m.URL()+"/test", strings.NewReader(`{"name": null, "ok": true}`))
		req.Header("Content-Type", "application/json")

		res, err := req.Do()

		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("array items", func(t *testing.T) {
		m := New()
		m.MustStart()

		defer m.Close()

		timestamp, err := time.Parse("2006-01-02T15:04:05.000Z", "2022-12-31T00:30:25.010Z")
		require.NoError(t, err)

		body := &jsonTestModel{
			Name:      "dev",
			OK:        true,
			Ordered:   []string{"dev", "qa", "devops"},
			Unordered: []string{"devops", "qa", "dev"},
			Diff:      []any{true, false, 10, 10.5, "dev", []string{"qa", "devops"}, map[string]any{"job": "none"}, nil},
			Level1: level1{
				Year:      2022,
				Timestamp: timestamp,
				Day:       0,
				Level2: level2{
					Value: 100.25,
				},
			},
		}

		scoped := m.MustMock(Postf("/test").
			Header("test", Equal("hello")).
			Body(
				Field("name", StrictEqual("dev")),
				Field("ok", StrictEqual(true)),
				Field("ordered", Equal([]any{"dev", "qa", "devops"})),
				Field("ordered", Contain("dev")),
				Field("ordered", Some(Equal("qa"))),
				Field("unordered", ItemsMatch([]any{"dev", "qa", "devops"})),
				Field("diff[0]", Truthy()),
				Field("diff[0]", Not(Falsy())),
				Field("diff[1]", Falsy()),
				Field("diff[1]", Not(Truthy())),
				Field("diff[1]", Not(IsNil())),
				Field("diff[2]", Equal(10)),
				Field("diff[2]", LessThan(11)),
				Field("diff[2]", LessThan(12)),
				Field("diff[3]", Equal(10.5)),
				Field("diff[4]", HasSuffix("v")),
				Field("diff[5]", HasLen(2)),
				Field("diff[5]", Equal([]any{"qa", "devops"})),
				Field("diff[6]", HasLen(1)),
				Field("diff[6]", HasKey("job")),
				Field("diff[7]", IsNil()),
				Field("level1.year", Equal(2022)),
				Field("level1.year", GreaterThan(2021)),
				Field("level1.year", GreaterThanOrEqual(2022)),
				Field("level1.level2.value", Equal(100.25)),
			).
			Reply(OK()))

		req := testutil.PostJSON(m.URL()+"/test", body)
		req.Header("test", "hello")

		res, err := req.Do()

		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestSimpleJSONValues(t *testing.T) {
	m := NewT(t)
	m.MustStart()

	testCases := []struct {
		name   string
		status int
		value  any
	}{
		{"string", http.StatusOK, "text"},
		{"bool", http.StatusCreated, true},
		{"int", http.StatusAccepted, 100},
		{"float32", http.StatusOK, float32(300.50)},
		{"float64", http.StatusPartialContent, 100.50},
		{"null", http.StatusAccepted, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scoped := m.MustMock(Postf("/test").Body(Equal(tc.value)).Reply(Status(tc.status)))

			res, err := testutil.PostJSON(m.URL()+"/test", tc.value).Do()

			require.NoError(t, err)
			require.Equal(t, tc.status, res.StatusCode)
			require.True(t, scoped.AssertCalled(t))
		})
	}
}

func TestMalformedJSON_ShouldMatchOtherFieldsAndContinue(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(Postf("/test").
		Header("test", StrictEqual("hello")).
		Reply(OK()))

	req := testutil.Post(m.URL()+"/test", strings.NewReader(`{"test": "malformed_json", "pass`))
	req.Header("test", "hello")
	req.Header("Content-Type", "application/json")

	res, err := req.Do()

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.True(t, scoped.HasBeenCalled())
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestJSONResponse(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	type payload struct {
		Language string `json:"language"`
		Active   bool   `json:"active"`
	}

	p := &payload{"go", true}

	m.MustMock(Getf("/test").Reply(OK().JSON(p)))

	res, err := testutil.Get(m.URL() + "/test").Do()
	require.NoError(t, err)

	defer res.Body.Close()

	var body payload

	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, p, &body)
}
