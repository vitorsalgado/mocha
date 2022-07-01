package mocha

import (
	"net/http"
	"net/url"

	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/mock"
)

type (
	MockBuilder struct {
		scenario              string
		scenarioRequiredState string
		scenarioNewState      string
		mock                  *mock.Mock
	}
)

const (
	weightNone    = 0
	weightVeryLow = 1
	weightLow     = 3
	weightRegular = 5
	weightHigh    = 7
)

func NewBuilder() *MockBuilder {
	return &MockBuilder{mock: mock.New()}
}

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
			Weight:      weightVeryLow,
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
			Weight:      weightRegular,
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
			Weight:      weightVeryLow,
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
			Weight:      weightLow,
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
			Weight:      weightVeryLow,
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
			Weight:      weightVeryLow,
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
				Weight:      weightHigh,
			})
	}

	return b
}

func (b *MockBuilder) FormField(field string, m matcher.Matcher[string]) *MockBuilder {
	b.mock.Expectations = append(b.mock.Expectations,
		mock.Expectation[string]{
			Name:        "form_field",
			ValuePicker: func(r *matcher.RequestInfo) string { return r.Request.Form.Get(field) },
			Matcher:     m,
			Weight:      weightVeryLow,
		})

	return b
}

func (b *MockBuilder) Form(m matcher.Matcher[url.Values]) *MockBuilder {
	b.mock.Expectations = append(b.mock.Expectations,
		mock.Expectation[url.Values]{
			Name:        "form",
			ValuePicker: func(r *matcher.RequestInfo) url.Values { return r.Request.Form },
			Matcher:     m,
			Weight:      weightLow,
		})

	return b
}

func (b *MockBuilder) Repeat(times int) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[any]{
			Name:        "repeat",
			ValuePicker: func(r *matcher.RequestInfo) any { return r.Request },
			Matcher:     matcher.Repeat[any](times),
			Weight:      weightNone,
		})

	return b
}

func (b *MockBuilder) Matches(m matcher.Matcher[any]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[any]{
			Name:        "any",
			ValuePicker: func(r *matcher.RequestInfo) any { return r.Request },
			Matcher:     m,
			Weight:      weightLow,
		})

	return b
}

func (b *MockBuilder) RequestMatches(m matcher.Matcher[*http.Request]) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		mock.Expectation[*http.Request]{
			Name:        "request",
			ValuePicker: func(r *matcher.RequestInfo) *http.Request { return r.Request },
			Matcher:     m,
			Weight:      weightLow,
		})

	return b
}

func (b *MockBuilder) StartScenario(scenario string) *MockBuilder {
	b.scenario = scenario
	b.scenarioRequiredState = matcher.ScenarioStarted
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

func (b *MockBuilder) PostAction(action mock.PostAction) *MockBuilder {
	b.mock.PostActions = append(b.mock.PostActions, action)
	return b
}

func (b *MockBuilder) Reply(reply mock.Reply) *MockBuilder {
	b.mock.Reply = reply
	return b
}

func (b *MockBuilder) Build() *mock.Mock {
	if b.scenario != "" {
		b.RequestMatches(matcher.Scenario[*http.Request](b.scenario, b.scenarioRequiredState, b.scenarioNewState))
	}

	return b.mock
}
