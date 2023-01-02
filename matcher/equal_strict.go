package matcher

import (
	"fmt"
	"reflect"
)

type equalStrictMatcher struct {
	expected any
}

func (m *equalStrictMatcher) Name() string {
	return "StrictEqual"
}

func (m *equalStrictMatcher) Match(v any) (*Result, error) {
	return &Result{
		Pass: reflect.DeepEqual(m.expected, v),
		Message: fmt.Sprintf("%s %s %v",
			hint(m.Name(), printExpected(m.expected)),
			_separator,
			printReceived(v)),
	}, nil
}

func (m *equalStrictMatcher) After() error {
	return nil
}

// StrictEqual returns true if matcher value and type are equal to the given parameter.
func StrictEqual(expected any) Matcher {
	return &equalStrictMatcher{expected: expected}
}

// StrictEqualf returns true if matcher value and type are equal to the given parameter.
// The expected value will be formatted with the provided format specifier.
func StrictEqualf(format string, a ...any) Matcher {
	return &equalStrictMatcher{expected: fmt.Sprintf(format, a...)}
}
