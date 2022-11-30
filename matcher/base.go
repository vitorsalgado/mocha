package matcher

import "net/http"

const (
	_separator = "=>"
)

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
	Match(value any) (Result, error)

	OnMockServed() error
}

type Result struct {
	OK              bool
	DescribeFailure func() string
}

type ComposableMatcher struct {
	M Matcher
}

func (m *ComposableMatcher) Name() string                { return m.M.Name() }
func (m *ComposableMatcher) Match(v any) (Result, error) { return m.M.Match(v) }
func (m *ComposableMatcher) OnMockServed() error         { return m.M.OnMockServed() }

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

func mismatch(failureMessageFunc func() string) Result {
	return Result{OK: false, DescribeFailure: failureMessageFunc}
}
