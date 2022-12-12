package mocha

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/reply"

	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestRace(t *testing.T) {
	m := newMock()
	jobs := 10
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(index int) {
			if index%2 == 0 {
				time.Sleep(100 * time.Millisecond)
			}

			m.Inc()
			wg.Done()
		}(i)

		m.Inc()
	}

	m.Inc()
	m.Inc()

	wg.Wait()

	assert.Equal(t, (jobs*2)+2, m.hits)
}

func TestMock(t *testing.T) {
	m := newMock()

	t.Run("should init enabled", func(t *testing.T) {
		assert.True(t, m.Enabled)
	})

	t.Run("should disable mock when calling .Disable()", func(t *testing.T) {
		m.Disable()
		assert.False(t, m.Enabled)

		m.Enable()
		assert.True(t, m.Enabled)
	})

	t.Run("should return called when it was hit", func(t *testing.T) {
		assert.False(t, m.HasBeenCalled())
		m.Inc()
		assert.True(t, m.HasBeenCalled())

		m.Dec()
		assert.False(t, m.HasBeenCalled())
	})
}

func TestMock_Matches(t *testing.T) {
	m := newMock()
	params := &RequestInfo{}

	cases := []struct {
		name     string
		value    any
		selector any
		expected bool
	}{
		{
			value:    "test",
			selector: "test",
			expected: true,
		},
		{
			value:    10.0,
			selector: 10.0,
			expected: true,
		},
		{
			value:    true,
			selector: true,
			expected: true,
		},
		{
			value:    map[string]any{"key": "value"},
			selector: map[string]any{"key": "value"},
			expected: true,
		},
		{
			value:    "test",
			selector: "dev",
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := m.requestMatches(params, []*expectation{{
				Matcher: Equal(tc.value),
				ValueSelector: func(r *RequestInfo) any {
					return tc.selector
				},
			}})
			assert.Equal(t, tc.expected, res.OK)
		})
	}

	t.Run("should return not matched and error when one of expectations returns error", func(t *testing.T) {
		// string
		res := m.requestMatches(params, []*expectation{{
			Matcher: Func(func(_ any) (bool, error) {
				return false, fmt.Errorf("fail")
			}),
			ValueSelector: func(r *RequestInfo) any {
				return "dev"
			},
		}})
		assert.False(t, res.OK)
	})

	t.Run("should return the sum of the matchers weight when it matches", func(t *testing.T) {
		// any
		res := m.requestMatches(params, []*expectation{
			{
				Matcher: Equal("test"),
				ValueSelector: func(r *RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: Equal("test"),
				ValueSelector: func(r *RequestInfo) any {
					return "test"
				},
				Weight: 1,
			},
			{
				Matcher: Equal(10.0),
				ValueSelector: func(r *RequestInfo) any {
					return 10.0
				},
				Weight: 2,
			},
		})
		assert.True(t, res.OK)
		assert.Equal(t, 5, res.Weight)
	})

	t.Run("should return the sum of the matchers weight when one of then doesnt matches", func(t *testing.T) {
		// any
		res := m.requestMatches(params, []*expectation{
			{
				Matcher: Equal("test"),
				ValueSelector: func(r *RequestInfo) any {
					return "test"
				},
				Weight: 2,
			},
			{
				Matcher: Equal("test"),
				ValueSelector: func(r *RequestInfo) any {
					return "dev"
				},
				Weight: 1,
			},
			{
				Matcher: Equal(10.0),
				ValueSelector: func(r *RequestInfo) any {
					return 10.0
				},
				Weight: 2,
			},
		})
		assert.False(t, res.OK)
		assert.Equal(t, 4, res.Weight)
	})
}

func TestMock_Build(t *testing.T) {
	m := newMock()
	m.Inc()
	m.Disable()

	mm, err := m.Build()

	assert.NoError(t, err)
	assert.Equal(t, m, mm)
}

func TestMock_MarshalJSON(t *testing.T) {
	file, err := os.Open(path.Join("testdata", "data.json"))
	assert.NoError(t, err)

	defer file.Close()

	// jzon := make(map[string]any)
	// err = json.NewDecoder(file).Decode(&jzon)
	// assert.NoError(t, err)

	m, err := Request().
		URLPath(Equal("/test")).
		Method(http.MethodPost).
		Query("q", EqualIgnoreCase("dev")).
		Query("sort", Equal("asc")).
		Header(header.ContentType, Contain("json")).
		Header(header.Accept, Contain("json")).
		Body(JSONPath("name", Equal("no-one"))).
		Body(JSONPath("active", Equal(true))).
		Times(5).
		Reply(reply.OK().
			BodyReader(file).
			Header("x-test", "ok").
			Header("x-dev", "nok").
			Cookie(&http.Cookie{Name: "hello", Value: "world"}).
			Cookie(&http.Cookie{Name: "hi", Value: "bye"})).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, m)

	b, err := m.MarshalJSON()

	assert.NoError(t, err)
	assert.NotNil(t, b)

	fmt.Println(string(b))
}
