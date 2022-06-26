package mocha

import (
	"github.com/vitorsalgado/mocha/mock"
	"net/http"
	"net/url"

	"github.com/vitorsalgado/mocha/matcher"
)

type (
	MockBuilder struct {
		scenario              string
		scenarioRequiredState string
		scenarioNewState      string
		mock                  *mock.Mock
	}
)

func NewBuilder() *MockBuilder                    { return &MockBuilder{mock: mock.New()} }
func Get(m matcher.Matcher[url.URL]) *MockBuilder { return NewBuilder().URL(m).Method(http.MethodGet) }
func Post(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodPost)
}
func Put(m matcher.Matcher[url.URL]) *MockBuilder { return NewBuilder().URL(m).Method(http.MethodPut) }
func Patch(u matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(u).Method(http.MethodPatch)
}
func Delete(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodDelete)
}
func Head(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodHead)
}
func Options(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodOptions)
}

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
		mock.Expectation[string]{
			Name:        "method",
			ValuePicker: func(r *matcher.RequestInfo) string { return r.Request.Method },
			Matcher:     matcher.EqualFold(method),
			Weight:      3,
		})

	return b
}

func (b *MockBuilder) URL(m matcher.Matcher[url.URL]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[url.URL]{
			Name:        "url",
			ValuePicker: func(r *matcher.RequestInfo) url.URL { return *r.Request.URL },
			Matcher:     m,
			Weight:      10,
		})

	return b
}

func (b *MockBuilder) Header(key string, m matcher.Matcher[string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[string]{
			Name:        "header",
			ValuePicker: func(r *matcher.RequestInfo) string { return r.Request.Header.Get(key) },
			Matcher:     m,
			Weight:      1,
		})

	return b
}

func (b *MockBuilder) Headers(m matcher.Matcher[map[string][]string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[map[string][]string]{
			Name:        "headers",
			ValuePicker: func(r *matcher.RequestInfo) map[string][]string { return r.Request.Header },
			Matcher:     m,
			Weight:      3,
		})

	return b
}

func (b *MockBuilder) Query(key string, m matcher.Matcher[string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[string]{
			Name:        "query",
			ValuePicker: func(r *matcher.RequestInfo) string { return r.Request.URL.Query().Get(key) },
			Matcher:     m,
			Weight:      1,
		})

	return b
}

func (b *MockBuilder) Queries(m matcher.Matcher[map[string][]string]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[map[string][]string]{
			Name:        "queries",
			ValuePicker: func(r *matcher.RequestInfo) map[string][]string { return r.Request.URL.Query() },
			Matcher:     m,
			Weight:      3,
		})

	return b
}

func (b *MockBuilder) Body(matchers ...matcher.Matcher[any]) *MockBuilder {
	for _, m := range matchers {
		b.mock.Expectations = append(b.mock.Expectations,
			mock.Expectation[any]{
				Name:        "body",
				ValuePicker: func(r *matcher.RequestInfo) any { return r.ParsedBody },
				Matcher:     m,
				Weight:      7,
			})
	}

	return b
}

func (b *MockBuilder) FormField(field string, m matcher.Matcher[string]) *MockBuilder {
	b.mock.Expectations = append(b.mock.Expectations,
		mock.Expectation[string]{
			Name:        "form",
			ValuePicker: func(r *matcher.RequestInfo) string { return r.Request.Form.Get(field) },
			Matcher:     m,
			Weight:      1,
		})

	return b
}

func (b *MockBuilder) Expect(m matcher.Matcher[*http.Request]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[*http.Request]{
			Name:        "request",
			ValuePicker: func(r *matcher.RequestInfo) *http.Request { return r.Request },
			Matcher:     m,
			Weight:      3,
		})

	return b
}

func (b *MockBuilder) StartScenario(scenario string) *MockBuilder {
	b.scenario = scenario
	b.scenarioRequiredState = ScenarioStarted
	return b
}

func (b *MockBuilder) ScenarioIs(scenario string) *MockBuilder {
	b.scenario = scenario
	return b
}

func (b *MockBuilder) ScenarioStateIs(requiredState string) *MockBuilder {
	b.scenarioRequiredState = requiredState
	return b
}

func (b *MockBuilder) ScenarioStateWillBe(newState string) *MockBuilder {
	b.scenarioNewState = newState
	return b
}

func (b *MockBuilder) Reply(reply mock.Reply) *MockBuilder {
	b.mock.Reply = reply
	return b
}

func (b *MockBuilder) Build() *mock.Mock {
	if b.scenario != "" {
		b.Expect(scenarioMatcher[*http.Request](b.scenario, b.scenarioRequiredState, b.scenarioNewState))
	}

	return b.mock
}
