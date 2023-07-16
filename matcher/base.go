package matcher

import (
	"errors"
	"strings"
)

// Matcher matches values.
type Matcher interface {
	Match(v any) (Result, error)
}

// OnAfterMockServed describes a Matcher that has post processes that need to be executed.
// AfterMockServed() function will be called after the mock HTTP response.
// Useful for stateful Matchers.
type OnAfterMockServed interface {
	AfterMockServed() error
}

// Result represents a Matcher expected.
type Result struct {
	// Pass defines if Matcher passed or not.
	Pass bool

	// Message describes why the associated Matcher did not pass.
	Message string
}

func runAfterMockServed(matchers ...Matcher) error {
	var errs []string

	for _, matcher := range matchers {
		m, ok := matcher.(OnAfterMockServed)
		if !ok {
			continue
		}

		err := m.AfterMockServed()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func success() Result                { return Result{Pass: true, Message: ""} }
func mismatch(message string) Result { return Result{Pass: false, Message: message} }
