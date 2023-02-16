package mocha

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/v3/matcher"
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

// AnyMethod creates a new empty Builder.
func AnyMethod() *MockBuilder {
	b := &MockBuilder{mock: newMock()}
	return b.MethodMatches(matcher.Anything())
}

// Get initializes a mock for GET method.
func Get(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodGet)
}

// Getf initializes a mock for GET method.
func Getf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodGet)
}

// Post initializes a mock for Post method.
func Post(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPost)
}

// Postf initializes a mock for Post method.
func Postf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPost)
}

// Put inits a mock for Put method.
func Put(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodPut)
}

// Putf initializes a mock for Put method.
func Putf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPut)
}

// Patch initializes a mock for Patch method.
func Patch(u matcher.Matcher) *MockBuilder {
	return Request().URL(u).Method(http.MethodPatch)
}

// Patchf initializes a mock for Patch method.
func Patchf(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodPatch)
}

// Delete initializes a mock for Delete method.
func Delete(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodDelete)
}

// Deletef initializes a mock for Delete method.
func Deletef(path string, a ...any) *MockBuilder {
	return Request().URLPathf(path, a...).Method(http.MethodDelete)
}

// Head initializes a mock for Head method.
func Head(m matcher.Matcher) *MockBuilder {
	return Request().URL(m).Method(http.MethodHead)
}

// Headf initializes a mock for Head method.
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
		Target:        _targetScheme,
		Key:           scheme,
		Matcher:       matcher.EqualIgnoreCase(scheme),
		ValueSelector: selectScheme,
		Weight:        _weightVeryLow,
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
		ValueSelector: selectMethod,
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
		ValueSelector: selectMethod,
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
		ValueSelector: selectURL,
		Matcher:       m,
		Weight:        _weightRegular,
	})

	return b
}

// URLf sets a matcher to the http.Request url.URL that compares the http.Request url.URL with given value.
// The expected value will be formatted with the provided format specifier.
func (b *MockBuilder) URLf(format string, a ...any) *MockBuilder {
	return b.URL(matcher.StrictEqual(fmt.Sprintf(format, a...)))
}

// URLPath defines a matcher to be applied to the url.URL path.
func (b *MockBuilder) URLPath(m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetURL,
		Key:           "url_path",
		ValueSelector: selectURLPath,
		Matcher:       m,
		Weight:        _weightRegular,
	})

	return b
}

// URLPathf sets a Matcher that compares the http.Request url.URL path with given value, ignoring case.
// The expected value will be formatted with the provided format specifier.
func (b *MockBuilder) URLPathf(format string, a ...any) *MockBuilder {
	return b.URLPath(matcher.StrictEqual(fmt.Sprintf(format, a...)))
}

// Header adds a matcher to a specific http.Header key.
func (b *MockBuilder) Header(key string, m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetHeader,
		Key:           key,
		ValueSelector: selectHeader(key),
		Matcher:       m,
		Weight:        _weightLow,
	})

	return b
}

// Headerf adds a matcher to a specific http.Header key.
func (b *MockBuilder) Headerf(key string, value string, a ...any) *MockBuilder {
	return b.Header(key, matcher.StrictEqual(fmt.Sprintf(value, a...)))
}

// Query defines a matcher to a specific query.
func (b *MockBuilder) Query(key string, m matcher.Matcher) *MockBuilder {
	b.appendExpectation(&expectation{
		Target:        _targetQuery,
		Key:           key,
		ValueSelector: selectQuery(key),
		Matcher:       m,
		Weight:        _weightVeryLow,
	})

	return b
}

// Queryf defines a matcher to a specific query.
func (b *MockBuilder) Queryf(key string, value string, a ...any) *MockBuilder {
	return b.Query(key, matcher.StrictEqual(fmt.Sprintf(value, a...)))
}

// Body adds matchers to the request body.
// If request contains a JSON body, you can provide multiple matchers to several fields.
// Example:
//
//	m.Body(JSONPath("name", EqualTo("test")), JSONPath("address.street", ToContains("nowhere")))
func (b *MockBuilder) Body(matcherList ...matcher.Matcher) *MockBuilder {
	var m matcher.Matcher
	if len(matcherList) == 0 {
		panic(".Body() func requires at least one matcher.Matcher")
	} else if len(matcherList) == 1 {
		m = matcherList[0]
	} else {
		m = matcher.AllOf(matcherList...)
	}

	b.appendExpectation(&expectation{
		Target:        _targetBody,
		ValueSelector: selectBody,
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
		ValueSelector: selectFormField(field),
		Matcher:       m,
		Weight:        _weightVeryLow,
	})

	return b
}

// FormFieldf defines a matcher for a specific form field by its key.
func (b *MockBuilder) FormFieldf(field string, value string, a ...any) *MockBuilder {
	return b.FormField(field, matcher.StrictEqual(fmt.Sprintf(value, a...)))
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
		ValueSelector: selectRawRequest,
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
func (b *MockBuilder) Reply(rep Reply) *MockBuilder {
	b.mock.Reply = rep
	return b
}

// Enable define if the Mock will enabled or disabled.
// All mocks are enabled by default.
func (b *MockBuilder) Enable(enabled bool) *MockBuilder {
	b.mock.Enabled = enabled
	return b
}

// Build builds a Mock with previously configured parameters.
// Used internally by Mocha.
func (b *MockBuilder) Build(_ *Mocha) (*Mock, error) {
	if len(b.mock.expectations) == 0 {
		return nil, fmt.Errorf("at least 1 request matcher must be set")
	}

	if b.mock.Reply == nil {
		return nil,
			fmt.Errorf("no reply set. use .Reply() or any equivalent to set the expected mock response")
	}

	if r, ok := b.mock.Reply.(replyValidation); ok {
		err := r.Validate()
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

// --
// Request Values Selectors
// --

func selectScheme(r *valueSelectorInput) any  { return r.URL.Scheme }
func selectMethod(r *valueSelectorInput) any  { return r.RawRequest.Method }
func selectURL(r *valueSelectorInput) any     { return r.URL }
func selectURLPath(r *valueSelectorInput) any { return r.URL.Path }
func selectHeader(k string) valueSelector {
	return func(r *valueSelectorInput) any { return r.RawRequest.Header.Get(k) }
}
func selectQuery(k string) valueSelector {
	return func(r *valueSelectorInput) any { return r.RawRequest.URL.Query().Get(k) }
}
func selectBody(r *valueSelectorInput) any { return r.ParsedBody }
func selectFormField(k string) valueSelector {
	return func(r *valueSelectorInput) any { return r.RawRequest.Form.Get(k) }
}
func selectRawRequest(r *valueSelectorInput) any { return r.RawRequest }
