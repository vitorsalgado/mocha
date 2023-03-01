package matcher

// Matcher defines an HTTP request matcher.
type Matcher interface {
	// Name names the Matcher.
	// Useful to give more context on non-matched requests.
	Name() string

	// Match is the function that does the actual matching logic.
	Match(value any) (*Result, error)
}

// OnAfterMockServed describes a Matcher that has post processes that need to be executed.
// AfterMockServed() function will be called after the mock HTTP response.
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

	// Ext defines extra information that gives more context to non-matched results.
	Ext []string
}
