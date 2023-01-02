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

	// After runs everytime the Mock that holds this Matcher is served.
	// Useful for stateful Matchers.
	After() error
}

// After describes a Matcher that has post processes that needs to be executed.
// The After() function will be called after mock HTTP response.
// Useful for stateful Matchers.
type After interface {
	After() error
}

// Result represents a Matcher result.
type Result struct {
	// Pass defines if Matcher passed or not.
	Pass bool

	// Message describes why the associated Matcher did not pass.
	Message string
}
