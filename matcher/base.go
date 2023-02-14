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
}

// OnAfterMockServed describes a Matcher that has post processes that needs to be executed.
// AfterMockServed() function will be called after mock HTTP response.
// Useful for stateful Matchers.
type OnAfterMockServed interface {
	AfterMockServed() error
}

// Result represents a Matcher result.
type Result struct {
	// Pass defines if Matcher passed or not.
	Pass bool

	// Message describes why the associated Matcher did not pass.
	Message string

	// Ext ...
	Ext []string
}
