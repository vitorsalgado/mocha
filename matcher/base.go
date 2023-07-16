package matcher

// Matcher matches values.
type Matcher interface {
	Match(v any) (Result, error)
}

// Result represents a Matcher expected.
type Result struct {
	Pass    bool   // Pass defines if Matcher passed or not.
	Message string // Message describes why the associated Matcher did not pass.
}

// OnMockSent describes a Matcher that has post processes that need to be executed after the matched mock is served.
// Dedicated for stateful matchers.
type OnMockSent interface {
	OnMockSent() error
}

func success() Result                { return Result{Pass: true, Message: ""} }
func mismatch(message string) Result { return Result{Pass: false, Message: message} }
