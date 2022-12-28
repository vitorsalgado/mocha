package mocha

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
	"github.com/vitorsalgado/mocha/v3/types"
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

// Getf inits a mock for GET method.
func Getf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodGet)
}

// Post inits a mock for Post method.
func Post(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPost)
}

// Postf inits a mock for Post method.
func Postf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPost)
}

// Put inits a mock for Put method.
func Put(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPut)
}

// Putf inits a mock for Put method.
func Putf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPut)
}

// Patch inits a mock for Patch method.
func Patch(u matcher.Matcher) *MockBuilder {
	return Request().URL(u).Method(http.MethodPatch)
}

// Patchf inits a mock for Patch method.
func Patchf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPatch)
}

// Delete inits a mock for Delete method.
func Delete(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodDelete)
}

// Deletef inits a mock for Delete method.
func Deletef(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodDelete)
}

// Head inits a mock for Head method.
func Head(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodHead)
}

// Headf inits a mock for Head method.
func Headf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodHead)
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

// Scheme sets a matcher.Matcher for the URL scheme part.
func (b *MockBuilder) Scheme(scheme string) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:  _targetScheme,
		Key:     scheme,
		Matcher: matcher.EqualIgnoreCase(scheme),
		ValueSelector: func(r *types.RequestValues) any {
			return r.URL.Scheme
		},
		Weight: _weightVeryLow,
	})

	return b
}

// Method sets the HTTP request method to be matched.
func (b *MockBuilder) Method(methods ...string) *MockBuilder {
	var m matcher.Matcher
	if len(methods) == 0 {
		panic(".Method() requires at least one HTTP Method")
	} else if len(methods) == 1 {
		m = matcher.EqualIgnoreCase(methods[0])
	} else {
		m = matcher.Some(methods)
	}

	b.appendExpectation(&expectation{
		Target:        _targetMethod,
		ValueSelector: func(r *types.RequestValues) any { return r.RawRequest.Method },
		Matcher:       m,
		Weight:        _weightNone,
	})

	return b
}

// MethodMatches defines a matcher.Matcher for the request method.
// Useful to set a Mock for multiple HTTP Request methods.
func (b *MockBuilder) MethodMatches(m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetMethod,
		ValueSelector: func(r *types.RequestValues) any { return r.RawRequest.Method },
		Matcher:       m,
		Weight:        _weightNone,
	})

	return b
}

// URL defines a matcher to be applied to the http.Request url.URL.
func (b *MockBuilder) URL(m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetURL,
		Key:           "url",
		ValueSelector: func(r *types.RequestValues) any { return r.RawRequest.URL },
		Matcher:       m,
		Weight:        _weightRegular,
	})

	return b
}

// URLf sets a matcher to the http.Request url.URL that compares the http.Request url.URL with given value.
// The expected value will be formatted with the provided format specifier.
func (b *MockBuilder) URLf(format string, a ...any) *MockBuilder {
	return b.URL(matcher.Equal(fmt.Sprintf(format, a...)))
}

// URLPath defines a matcher to be applied to the url.URL path.
func (b *MockBuilder) URLPath(m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetURL,
		Key:           "url_path",
		ValueSelector: func(r *types.RequestValues) any { return r.RawRequest.URL.Path },
		Matcher:       m,
		Weight:        _weightRegular,
	})

	return b
}

// URLPathf sets a Matcher that compares the http.Request url.URL path with given value, ignoring case.
// The expected value will be formatted with the provided format specifier.
func (b *MockBuilder) URLPathf(format string, a ...any) *MockBuilder {
	return b.URLPath(matcher.Equal(fmt.Sprintf(format, a...)))
}

// Header adds a matcher to a specific http.Header key.
func (b *MockBuilder) Header(key string, m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetHeader,
		Key:           key,
		ValueSelector: func(r *types.RequestValues) any { return r.RawRequest.Header.Get(key) },
		Matcher:       m,
		Weight:        _weightLow,
	})

	return b
}

// Headerf adds a matcher to a specific http.Header key.
func (b *MockBuilder) Headerf(key string, value string, a ...any) *MockBuilder {
	return b.Header(key, matcher.Equal(fmt.Sprintf(value, a...)))
}

// Query defines a matcher to a specific query.
func (b *MockBuilder) Query(key string, m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetQuery,
		Key:           key,
		ValueSelector: func(r *types.RequestValues) any { return r.RawRequest.URL.Query().Get(key) },
		Matcher:       m,
		Weight:        _weightVeryLow,
	})

	return b
}

// Queryf defines a matcher to a specific query.
func (b *MockBuilder) Queryf(key string, value string, a ...any) *MockBuilder {
	return b.Query(key, matcher.Equal(fmt.Sprintf(value, a...)))
}

// Body adds matchers to the request body.
// If request contains a JSON body, you can provide multiple matchers to several fields.
// Example:
//
//	m.ParsedBody(JSONPath("name", EqualTo("test")), JSONPath("address.street", ToContains("nowhere")))
func (b *MockBuilder) Body(matcherList ...matcher.Matcher) *MockBuilder {
	var m matcher.Matcher
	if len(matcherList) == 0 {
		panic(".ParsedBody() func requires at least one matcher.Matcher")
	} else if len(matcherList) == 1 {
		m = matcherList[0]
	} else {
		m = matcher.AllOf(matcherList...)
	}

	b.appendExpectation(&expectation{
		Target:        _targetBody,
		ValueSelector: func(r *types.RequestValues) any { return r.ParsedBody },
		Matcher:       m,
		Weight:        _weightHigh,
	})

	return b
}

// FormField defines a matcher for a specific form field by its key.
func (b *MockBuilder) FormField(field string, m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetForm,
		Key:           field,
		ValueSelector: func(r *types.RequestValues) any { return r.RawRequest.Form.Get(field) },
		Matcher:       m,
		Weight:        _weightVeryLow,
	})

	return b
}

// FormFieldf defines a matcher for a specific form field by its key.
func (b *MockBuilder) FormFieldf(field string, value string, a ...any) *MockBuilder {
	return b.FormField(field, matcher.Equal(fmt.Sprintf(value, a...)))
}

// Times defines to total times that a mock should be served, if request matches.
func (b *MockBuilder) Times(times int) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:  _targetRequest,
		Matcher: matcher.Repeat(times),
		Weight:  _weightNone,
	})
	return b
}

// RequestMatches defines matcher.Matcher to be applied to a http.Request.
func (b *MockBuilder) RequestMatches(m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetRequest,
		ValueSelector: func(r *types.RequestValues) any { return r.ParsedBody },
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

// Map adds a Mapper that allows modifying the response after it was built.
// Multiple mappers can be added.
// Map doesn't work with reply.From or Proxy.
func (b *MockBuilder) Map(mapper Mapper) *MockBuilder {
	b.mock.Mappers = append(b.mock.Mappers, mapper)
	return b
}

// Reply defines a response mock to be served if this mock matches to a request.
func (b *MockBuilder) Reply(rep reply.Reply) *MockBuilder {
	b.mock.Reply = rep
	return b
}

// Enabled define if the Mock will enabled or disabled.
// All mocks are enabled by default.
func (b *MockBuilder) Enabled(enabled bool) *MockBuilder {
	b.mock.Enabled = enabled
	return b
}

// Build builds a Mock with previously configured parameters.
// Used internally by Mocha.
func (b *MockBuilder) Build() (*Mock, error) {
	if len(b.mock.expectations) == 0 {
		return nil, fmt.Errorf("at least 1 request matcher must be set")
	}

	if b.mock.Reply == nil {
		return nil,
			fmt.Errorf("no reply set. use .Reply() or any equivalent to set the expected mock response")
	}

	if r, ok := b.mock.Reply.(reply.Pre); ok {
		err := r.Pre()
		if err != nil {
			return nil, err
		}
	}

	if b.scenario != "" {
		b.appendExpectation(&expectation{
			Target:  _targetRequest,
			Matcher: matcher.Scenario(b.scenario, b.scenarioRequiredState, b.scenarioNewState),
		})
	}

	return b.mock, nil
}

func (b *MockBuilder) appendExpectation(e *expectation) {
	b.mock.expectations = append(b.mock.expectations, e)
}
