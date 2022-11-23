package expect

import "net/http"

// RequestInfo implements HTTP request information to be passed to each Matcher.
type RequestInfo struct {
	// Request is the actual http.Request.
	Request *http.Request

	// ParsedBody is http.Request parsed body.
	// Value of parsed body can vary depending on the mocha.RequestBodyParser that parsed the request.
	ParsedBody any
}

// Matcher defines request matchers.
// Request matchers are used to match requests in order to find a mock to serve a stub response.
type Matcher interface {
	Name() string

	// Match is the function that does the actual matching logic.
	Match(value any) (bool, error)

	// DescribeFailure gives more context of why the Matcher failed to match a given value.
	DescribeFailure(value any) string

	OnMockServed() error
}

type ComposableMatcher struct {
	M Matcher
}

func (m *ComposableMatcher) Name() string              { return m.M.Name() }
func (m *ComposableMatcher) Match(v any) (bool, error) { return m.M.Match(v) }
func (m *ComposableMatcher) OnMockServed() error       { return nil }
func (m *ComposableMatcher) DescribeFailure(v any) string {
	return m.M.DescribeFailure(v)
}

// And compose the current Matcher with another one using the "and" operator.
func (m *ComposableMatcher) And(and Matcher) *ComposableMatcher {
	return Compose(AllOf(m, and))
}

// Or compose the current Matcher with another one using the "or" operator.
func (m *ComposableMatcher) Or(or Matcher) *ComposableMatcher {
	return Compose(AnyOf(m, or))
}

// Xor compose the current Matcher with another one using the "xor" operator.
func (m *ComposableMatcher) Xor(and Matcher) *ComposableMatcher {
	return Compose(XOR(m, and))
}

func Compose(base Matcher) *ComposableMatcher {
	return &ComposableMatcher{M: base}
}
