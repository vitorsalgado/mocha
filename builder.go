package mocha

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/params"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type Deps struct {
	ScenarioStore expect.ScenarioStorage
}

// MockBuilder is a builder for mock.Mock.
type MockBuilder struct {
	mock                  *Mock
	scenario              string
	scenarioNewState      string
	scenarioRequiredState string
}

// Request creates a new empty MockBuilder.
func Request() *MockBuilder {
	return &MockBuilder{mock: newMock()}
}

// Get inits a mock for GET method.
func Get(m expect.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodGet)
}

// Post inits a mock for Post method.
func Post(m expect.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPost)
}

// Put inits a mock for Put method.
func Put(m expect.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPut)
}

// Patch inits a mock for Patch method.
func Patch(u expect.Matcher) *MockBuilder {
	return Request().URL(u).Method(http.MethodPatch)
}

// Delete inits a mock for Delete method.
func Delete(m expect.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodDelete)
}

// Head inits a mock for Head method.
func Head(m expect.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodHead)
}

// Options inits a mock for Options method.
func Options(m expect.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodOptions)
}

// Name defines a name for the mock.
// Useful to debug.
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
		Expectation{
			Target:        "method",
			ValueSelector: func(r *expect.RequestInfo) any { return r.Request.Method },
			Matcher:       expect.ToEqualFold(method),
			Weight:        _weightNone,
		})

	return b
}

// URL defines a matcher to be applied to the http.Request url.URL.
func (b *MockBuilder) URL(m expect.Matcher) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation{
			Target:        "url",
			ValueSelector: func(r *expect.RequestInfo) any { return r.Request.URL },
			Matcher:       m,
			Weight:        _weightRegular,
		})

	return b
}

// Header adds a matcher to a specific http.Header key.
func (b *MockBuilder) Header(key string, m expect.Matcher) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation{
			Target:        "header",
			ValueSelector: func(r *expect.RequestInfo) any { return r.Request.Header.Get(key) },
			Matcher:       m,
			Weight:        _weightLow,
		})

	return b
}

// Query defines a matcher to a specific query.
func (b *MockBuilder) Query(key string, m expect.Matcher) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation{
			Target:        "query",
			ValueSelector: func(r *expect.RequestInfo) any { return r.Request.URL.Query().Get(key) },
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
func (b *MockBuilder) Body(matcherList ...expect.Matcher) *MockBuilder {
	for _, m := range matcherList {
		b.mock.Expectations = append(b.mock.Expectations,
			Expectation{
				Target:        "body",
				ValueSelector: func(r *expect.RequestInfo) any { return r.ParsedBody },
				Matcher:       m,
				Weight:        _weightHigh,
			})
	}

	return b
}

// FormField defines a matcher for a specific form field by its key.
func (b *MockBuilder) FormField(field string, m expect.Matcher) *MockBuilder {
	b.mock.Expectations = append(b.mock.Expectations,
		Expectation{
			Target:        "form",
			ValueSelector: func(r *expect.RequestInfo) any { return r.Request.Form.Get(field) },
			Matcher:       m,
			Weight:        _weightVeryLow,
		})

	return b
}

// Repeat defines to total times that a mock should be served, if request matches.
func (b *MockBuilder) Repeat(times int64) *MockBuilder {
	b.mock.Expectations = append(b.mock.Expectations,
		Expectation{
			Target:        "request",
			ValueSelector: func(r *expect.RequestInfo) any { return r.Request },
			Matcher:       expect.Repeat(times),
			Weight:        _weightNone,
		})
	return b
}

// RequestMatches defines expect.Matcher to be applied to a http.Request.
func (b *MockBuilder) RequestMatches(m expect.Matcher) *MockBuilder {
	b.mock.Expectations = append(
		b.mock.Expectations,
		Expectation{
			Target:        "request",
			ValueSelector: func(r *expect.RequestInfo) any { return r.Request },
			Matcher:       m,
			Weight:        _weightLow,
		})

	return b
}

// StartScenario sets that this mock will start a new scenario with the given name.
func (b *MockBuilder) StartScenario(name string) *MockBuilder {
	b.scenario = name
	b.scenarioRequiredState = expect.ScenarioStateStarted
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

// Reply defines a response mock to be served if this mock matches to a request.
func (b *MockBuilder) Reply(rep reply.Reply) *MockBuilder {
	b.mock.Reply = rep
	return b
}

// ReplyFunction defines a function to will build the response mock.
func (b *MockBuilder) ReplyFunction(fn func(*http.Request, reply.M, params.P) (*reply.Response, error)) *MockBuilder {
	b.mock.Reply = reply.Function(fn)
	return b
}

// ReplyJust sets the mock to return a simple response with the given status code.
// Optionally, you can provide a reply as well. The status provided in the first parameter will prevail.
// Only the first reply will be used.
func (b *MockBuilder) ReplyJust(status int, r ...*reply.StdReply) *MockBuilder {
	if len(r) > 0 {
		rep := r[0]
		rep.Status(status)

		b.mock.Reply = rep
	} else {
		b.mock.Reply = reply.Status(status)
	}

	return b
}

// Build builds a mock.Mock with previously configured parameters.
// Used internally by Mocha.
func (b *MockBuilder) Build(deps *Deps) *Mock {
	if b.scenario != "" {
		b.mock.Expectations = append(b.mock.Expectations,
			Expectation{
				Target: "scenario",
				ValueSelector: func(r *expect.RequestInfo) any {
					return r.Request
				},
				Matcher: expect.
					Scenario(deps.ScenarioStore)(
					b.scenario,
					b.scenarioRequiredState,
					b.scenarioNewState),
			})
	}

	return b.mock
}
