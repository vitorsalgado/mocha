package mocha

import (
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

var _ Builder = (*MockBuilder)(nil)

// MockBuilder is a builder for Mock.
type MockBuilder struct {
	mock                  *Mock
	scenario              string
	scenarioNewState      string
	scenarioRequiredState string
}

// Request creates a new empty Builder.
func Request() *MockBuilder {
	return &MockBuilder{mock: newMock()}
}

// Get inits a mock for GET method.
func Get(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodGet)
}

// Post inits a mock for Post method.
func Post(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPost)
}

// Put inits a mock for Put method.
func Put(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPut)
}

// Patch inits a mock for Patch method.
func Patch(u matcher.Matcher) *MockBuilder {
	return Request().URL(u).Method(http.MethodPatch)
}

// Delete inits a mock for Delete method.
func Delete(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodDelete)
}

// Head inits a mock for Head method.
func Head(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodHead)
}

// Options inits a mock for Options method.
func Options(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodOptions)
}

// Name defines a name for the mock.
// Useful for debug.
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
	b.mock.expectations = append(
		b.mock.expectations,
		&expectation{
			Target:        "method",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request.Method },
			Matcher:       matcher.EqualIgnoreCase(method),
			Weight:        _weightNone,
		})

	return b
}

// URL defines a matcher to be applied to the http.Request url.URL.
func (b *MockBuilder) URL(m matcher.Matcher) *MockBuilder {
	b.mock.expectations = append(
		b.mock.expectations,
		&expectation{
			Target:        "url",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request.URL },
			Matcher:       m,
			Weight:        _weightRegular,
		})

	return b
}

// URLPath defines a matcher to be applied to the url.URL path.
func (b *MockBuilder) URLPath(m matcher.Matcher) *MockBuilder {
	b.mock.expectations = append(
		b.mock.expectations,
		&expectation{
			Target:        "url",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request.URL.Path },
			Matcher:       m,
			Weight:        _weightRegular,
		})

	return b
}

// Header adds a matcher to a specific http.Header key.
func (b *MockBuilder) Header(key string, m matcher.Matcher) *MockBuilder {
	b.mock.expectations = append(
		b.mock.expectations,
		&expectation{
			Target:        "header",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request.Header.Get(key) },
			Matcher:       m,
			Weight:        _weightLow,
		})

	return b
}

// Query defines a matcher to a specific query.
func (b *MockBuilder) Query(key string, m matcher.Matcher) *MockBuilder {
	b.mock.expectations = append(
		b.mock.expectations,
		&expectation{
			Target:        "query",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request.URL.Query().Get(key) },
			Matcher:       m,
			Weight:        _weightVeryLow,
		})

	return b
}

// Body adds matchers to the request body.
// If request contains a JSON body, you can provide multiple matchers to several fields.
// Example:
//
//	m.Body(JSONPath("name", EqualTo("test")), JSONPath("address.street", ToContains("nowhere")))
func (b *MockBuilder) Body(matcherList ...matcher.Matcher) *MockBuilder {
	for _, m := range matcherList {
		b.mock.expectations = append(b.mock.expectations,
			&expectation{
				Target:        "body",
				ValueSelector: func(r *matcher.RequestInfo) any { return r.ParsedBody },
				Matcher:       m,
				Weight:        _weightHigh,
			})
	}

	return b
}

// FormField defines a matcher for a specific form field by its key.
func (b *MockBuilder) FormField(field string, m matcher.Matcher) *MockBuilder {
	b.mock.expectations = append(b.mock.expectations,
		&expectation{
			Target:        "form",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request.Form.Get(field) },
			Matcher:       m,
			Weight:        _weightVeryLow,
		})

	return b
}

// Repeat defines to total times that a mock should be served, if request matches.
func (b *MockBuilder) Repeat(times int) *MockBuilder {
	b.mock.expectations = append(b.mock.expectations,
		&expectation{
			Target:        "request",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request },
			Matcher:       matcher.Repeat(times),
			Weight:        _weightNone,
		})
	return b
}

// RequestMatches defines matcher.Matcher to be applied to a http.Request.
func (b *MockBuilder) RequestMatches(m matcher.Matcher) *MockBuilder {
	b.mock.expectations = append(
		b.mock.expectations,
		&expectation{
			Target:        "request",
			ValueSelector: func(r *matcher.RequestInfo) any { return r.Request },
			Matcher:       m,
			Weight:        _weightLow,
		})

	return b
}

// StartScenario sets that this mock will start a new scenario with the given name.
func (b *MockBuilder) StartScenario(name string) *MockBuilder {
	b.scenario = name
	b.scenarioRequiredState = matcher.ScenarioStateStarted
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

// PostAction adds a post action to be executed after the mocked response is served.
func (b *MockBuilder) PostAction(action PostAction) *MockBuilder {
	b.mock.PostActions = append(b.mock.PostActions, action)
	return b
}

// Delay sets a delay time before serving the mocked response.
func (b *MockBuilder) Delay(duration time.Duration) *MockBuilder {
	b.mock.Delay = duration
	return b
}

// Reply defines a response mock to be served if this mock matches to a request.
func (b *MockBuilder) Reply(rep reply.Reply) *MockBuilder {
	b.mock.Reply = rep
	return b
}

// ReplyFunc defines a function to will build the response mock.
func (b *MockBuilder) ReplyFunc(
	fn func(http.ResponseWriter, *http.Request) (*reply.Response, error),
) *MockBuilder {
	b.mock.Reply = reply.Function(fn)
	return b
}

// ReplyJust sets the mock to return a simple response with the given status code.
// Optionally, you can provide a reply as well.
func (b *MockBuilder) ReplyJust(status int, r *reply.StdReply) *MockBuilder {
	if r == nil {
		b.mock.Reply = reply.Status(status)
	} else {
		b.mock.Reply = r.Status(status)
	}

	return b
}

// Build builds a Mock with previously configured parameters.
// Used internally by Mocha.
func (b *MockBuilder) Build() *Mock {
	if b.scenario != "" {
		b.mock.expectations = append(b.mock.expectations,
			&expectation{
				Target: "scenario",
				ValueSelector: func(r *matcher.RequestInfo) any {
					return r.Request
				},
				Matcher: matcher.Scenario(b.scenario, b.scenarioRequiredState, b.scenarioNewState),
			})
	}

	return b.mock
}
