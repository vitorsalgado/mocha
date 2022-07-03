package mocha

import (
	"net/http"
	"net/url"

	"github.com/vitorsalgado/mocha/internal/scenario"
	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/mock"
)

type (
	// MockBuilder is a builder for mock.Mock.
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

// Get inits a mock for GET method.
func Get(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodGet)
}

// Post inits a mock for Post method.
func Post(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodPost)
}

// Put inits a mock for Put method.
func Put(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodPut)
}

// Patch inits a mock for Patch method.
func Patch(u matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(u).Method(http.MethodPatch)
}

// Delete inits a mock for Delete method.
func Delete(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodDelete)
}

// Head inits a mock for Head method.
func Head(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodHead)
}

// Options inits a mock for Options method.
func Options(m matcher.Matcher[url.URL]) *MockBuilder {
	return NewBuilder().URL(m).Method(http.MethodOptions)
}

// Name defines a name for the mock.
// Useful to debug your mocks.
func (b *MockBuilder) Name(name string) *MockBuilder {
	b.mock.Name = name
	return b
}

// Priority sets the priority of the mock.
// A higher priority will take precedence during request matching.
func (b *MockBuilder) Priority(p int) *MockBuilder {
	b.mock.Priority = p
	return b
}

// Method sets the HTTP request method to be matched.
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

// URL defines a matcher to be applied to the http.Request url.URL.
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

// Header adds a matcher to a specific http.Header key.
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

// Headers defines a matcher to be applied against all request headers.
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

// Query defines a matcher to a specific query.
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

// Queries defines a matcher to be applied to all request queries.
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

// Body adds matchers to the request body.
// If request contains a JSON body, you can provide multiple matchers to several fields.
// Example:
//	m.Body(JSONPath("name", EqualTo("test")), JSONPath("address.street", Contains("nowhere")))
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

// FormField defines a matcher for a specific form field by its key.
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

// Form defines a matcher to the request form url.Values.
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

// Repeat defines to total times that a mock should be served, if request matches.
func (b *MockBuilder) Repeat(times int) *MockBuilder {
	b.mock.AfterExpectations = append(
		b.mock.AfterExpectations,
		mock.Expectation[any]{
			Name:        "repeat",
			ValuePicker: func(r *matcher.RequestInfo) any { return r.Request },
			Matcher:     matcher.Repeat(times),
			Weight:      weightNone,
		})

	return b
}

// Matches sets a custom matcher, not necessary tied to a specific request parameter.
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

// RequestMatches defines matcher.Matcher to be applied to a http.Request.
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

// StartScenario sets that this mock will start a new scenario with the given name.
func (b *MockBuilder) StartScenario(name string) *MockBuilder {
	b.scenario = name
	b.scenarioRequiredState = scenario.StateStarted
	return b
}

// ScenarioIs mark this mock to be used only within the given scenario.
func (b *MockBuilder) ScenarioIs(scenario string) *MockBuilder {
	b.scenario = scenario
	return b
}

// ScenarioStateIs mark this mock to be served only if the scenario state is equal to the given required state.
func (b *MockBuilder) ScenarioStateIs(requiredState string) *MockBuilder {
	b.scenarioRequiredState = requiredState
	return b
}

// ScenarioStateWillBe defines the state of the scenario after this mock is matched, making the scenario flow continue.
func (b *MockBuilder) ScenarioStateWillBe(newState string) *MockBuilder {
	b.scenarioNewState = newState
	return b
}

// MatchAfter adds a matcher.Matcher to be run after the standard matchers and before serving the mocked response.
// After matchers are mostly used in special cases, like when they need to keep data that should not be evaluated all the time.
func (b *MockBuilder) MatchAfter(m matcher.Matcher[any]) *MockBuilder {
	b.mock.AfterExpectations = append(
		b.mock.AfterExpectations,
		mock.Expectation[any]{
			Name:        "after",
			ValuePicker: func(r *matcher.RequestInfo) any { return r.Request },
			Matcher:     m,
			Weight:      weightNone,
		})

	return b
}

// PostAction adds a post action to be executed after the mocked response is served.
func (b *MockBuilder) PostAction(action mock.PostAction) *MockBuilder {
	b.mock.PostActions = append(b.mock.PostActions, action)
	return b
}

// Reply defines a response stub to be served if this mock matches to a request.
func (b *MockBuilder) Reply(reply mock.Reply) *MockBuilder {
	b.mock.Reply = reply
	return b
}

// Build builds a mock.Mock with previously configured parameters.
// Used internally by Mocha.
func (b *MockBuilder) Build() *mock.Mock {
	if b.scenario != "" {
		b.RequestMatches(scenario.Scenario[*http.Request](b.scenario, b.scenarioRequiredState, b.scenarioNewState))
	}

	return b.mock
}
