package matcher

import "github.com/vitorsalgado/mocha/v3/types"

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

	// AfterMockSent runs everytime the Mock that holds this Matcher is served.
	// Useful for stateful Matchers.
	AfterMockSent() error

	// Raw "serializes" the Matcher to its raw format: ["matcher name", ...<parameters (any)>]
	// The raw format is used to save, load and build mocks from external sources, like JSON, YAML etc.
	// First array item must be the Matcher unique name.
	// Second item onwards are the Matcher parameters and ca be anything, including another Matcher
	// that follows the same spec.
	Raw() types.RawValue
}

// Result represents a Matcher result.
type Result struct {
	// Pass defines if Matcher passed or not.
	Pass bool

	// Message is function that should return the failure description.
	Message func() string
}
