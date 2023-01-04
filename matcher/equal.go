package matcher

import (
	"fmt"
)

type equalMatcher struct {
	expected any
}

func (m *equalMatcher) Name() string {
	return "Equal"
}

func (m *equalMatcher) Match(v any) (*Result, error) {
	if equalValues(m.expected, v) {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: fmt.Sprintf(
			"%s %s %v",
			hint(m.Name(), printExpected(m.expected)),
			_separator,
			printReceived(v),
		)},
		nil
}

// Equal asserts that the given expectation is equal to the incoming request value.
// It considers equivalent value. Eg.: float64(10) is equal to int(10).
func Equal(expected any) Matcher {
	return &equalMatcher{expected: expected}
}

// Equalf returns true if matcher value is equal to the given parameter value.
// It considers equivalent value. Eg.: float64(10) is equal to int(10).
// The expected value will be formatted with the provided format specifier.
func Equalf(format string, a ...any) Matcher {
	return &equalMatcher{expected: fmt.Sprintf(format, a...)}
}
