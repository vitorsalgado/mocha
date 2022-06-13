package mocha

import (
	"net/http"
	"net/url"
)

type MockBuilder struct {
	mock *Mock
}

func NewBuilder() *MockBuilder               { return &MockBuilder{mock: NewMock()} }
func Get(m Matcher[url.URL]) *MockBuilder    { return NewBuilder().URL(m).Method(http.MethodGet) }
func Post(m Matcher[url.URL]) *MockBuilder   { return NewBuilder().URL(m).Method(http.MethodPost) }
func Put(m Matcher[url.URL]) *MockBuilder    { return NewBuilder().URL(m).Method(http.MethodPut) }
func Patch(u Matcher[url.URL]) *MockBuilder  { return NewBuilder().URL(u).Method(http.MethodPatch) }
func Delete(m Matcher[url.URL]) *MockBuilder { return NewBuilder().URL(m).Method(http.MethodDelete) }
func Head(m Matcher[url.URL]) *MockBuilder   { return NewBuilder().URL(m).Method(http.MethodHead) }

func (b *MockBuilder) Name(name string) *MockBuilder {
	b.mock.Name = name
	return b
}

func (b *MockBuilder) Priority(p int) *MockBuilder {
	b.mock.Priority = p
	return b
}

func (b *MockBuilder) Method(method string) *MockBuilder {
	var matcher Matcher[string]
	var weight = 0

	if method == "*" {
		matcher = Anything[string]()
	} else {
		matcher = EqualFold(method)
	}

	weight = 3
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[string]{
			Name:    "method",
			Pick:    func(r *http.Request) string { return r.Method },
			Matcher: matcher,
			Weight:  weight,
		})

	return b
}

func (b *MockBuilder) URL(matcher Matcher[url.URL]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[url.URL]{
			Name:    "url",
			Pick:    func(r *http.Request) url.URL { return *r.URL },
			Matcher: matcher,
			Weight:  10,
		})

	return b
}

func (b *MockBuilder) Header(key string, matcher Matcher[string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[string]{
			Name:    "header",
			Pick:    func(r *http.Request) string { return r.Header.Get(key) },
			Matcher: matcher,
			Weight:  1,
		})

	return b
}

func (b *MockBuilder) Headers(matcher Matcher[map[string][]string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[map[string][]string]{
			Name:    "headers",
			Pick:    func(r *http.Request) map[string][]string { return r.Header },
			Matcher: matcher,
			Weight:  3,
		})

	return b
}

func (b *MockBuilder) Query(key string, matcher Matcher[string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[string]{
			Name:    "query",
			Pick:    func(r *http.Request) string { return r.URL.Query().Get(key) },
			Matcher: matcher,
			Weight:  1,
		})

	return b
}

func (b *MockBuilder) Queries(matcher Matcher[map[string][]string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[map[string][]string]{
			Name:    "queries",
			Pick:    func(r *http.Request) map[string][]string { return r.URL.Query() },
			Matcher: matcher,
			Weight:  3,
		})

	return b
}

func (b *MockBuilder) Expect(matcher Matcher[*http.Request]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[*http.Request]{
			Name:    "request",
			Pick:    func(r *http.Request) *http.Request { return r },
			Matcher: matcher,
			Weight:  3,
		})

	return b
}

func (b *MockBuilder) Scenario(name, requiredState, newState string) *MockBuilder {
	b.Expect(scenarioMatcher[*http.Request](name, requiredState, newState))
	return b
}

func (b *MockBuilder) Reply(reply Reply) *MockBuilder {
	b.mock.ResFn = reply.Build()
	return b
}

func (b *MockBuilder) Build() *Mock {
	return b.mock
}
