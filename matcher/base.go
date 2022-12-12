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
	// Name names the Matcher.
	// Used to give more context on non-matched requests.
	Name() string

	// Match is the function that does the actual matching logic.
	Match(value any) (*Result, error)

	// OnMockServed runs everytime the Mock that holds this Matcher is served.
	// Useful for stateful Matchers.
	OnMockServed() error

	// Spec serializes the Matcher to the format: ["matcher name", ...<parameters (any)>]
	Spec() any
}

type Result struct {
	OK              bool
	DescribeFailure func() string
}

// Values stores Matcher values and provides means to access then.
type Values struct {
	V any
}

func (v Values) Interface() any {
	return v.V
}

func (v Values) String() string {
	return v.V.(string)
}
