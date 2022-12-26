package matcher

const (
	_separator = "=>"
)

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
}

// Result represents a Matcher result.
type Result struct {
	// Pass defines if Matcher passed or not.
	Pass bool

	// Message is function that should return the failure description.
	Message func() string
}
