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
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[string]{
			Name:    "method",
			Pick:    func(r *MockRequest) string { return r.RawRequest.Method },
			Matcher: EqualFold(method),
			Weight:  3,
		})

	return b
}

func (b *MockBuilder) URL(matcher Matcher[url.URL]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[url.URL]{
			Name:    "url",
			Pick:    func(r *MockRequest) url.URL { return *r.RawRequest.URL },
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
			Pick:    func(r *MockRequest) string { return r.RawRequest.Header.Get(key) },
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
			Pick:    func(r *MockRequest) map[string][]string { return r.RawRequest.Header },
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
			Pick:    func(r *MockRequest) string { return r.RawRequest.URL.Query().Get(key) },
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
			Pick:    func(r *MockRequest) map[string][]string { return r.RawRequest.URL.Query() },
			Matcher: matcher,
			Weight:  3,
		})

	return b
}

func (b *MockBuilder) Body(matchers ...Matcher[any]) *MockBuilder {
	for _, matcher := range matchers {
		b.mock.Expectations = append(b.mock.Expectations,
			Expectation[any]{
				Name:    "body",
				Pick:    func(r *MockRequest) any { return r.Body },
				Matcher: matcher,
				Weight:  7,
			})
	}

	return b
}

func (b *MockBuilder) FormField(field string, matcher Matcher[string]) *MockBuilder {
	b.mock.Expectations = append(b.mock.Expectations,
		Expectation[string]{
			Name:    "form",
			Pick:    func(r *MockRequest) string { return r.RawRequest.Form.Get(field) },
			Matcher: matcher,
			Weight:  1,
		})

	return b
}

func (b *MockBuilder) Expect(matcher Matcher[*http.Request]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation[*http.Request]{
			Name:    "request",
			Pick:    func(r *MockRequest) *http.Request { return r.RawRequest },
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
	b.mock.Responder = reply.Build()
	return b
}

func (b *MockBuilder) Build() *Mock {
	return b.mock
}
