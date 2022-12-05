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
	Match(value any) (*Result, error)

	OnMockServed() error
}

type Result struct {
	OK              bool
	DescribeFailure func() string
}
