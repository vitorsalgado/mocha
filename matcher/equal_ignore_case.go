package matcher

import (
	"fmt"
	"strings"
)

type equalIgnoreCaseMatcher struct {
	expected string
}

func (m *equalIgnoreCaseMatcher) Match(v any) (Result, error) {
	txt := ""
	if v != nil {
		var ok bool
		txt, ok = v.(string)
		if !ok {
			return Result{}, fmt.Errorf("eqi: value must be a string. got: %T", v)
		}
	}

	if strings.EqualFold(m.expected, txt) {
		return Result{Pass: true}, nil
	}

	return mismatch(strings.Join([]string{"Eqi(", m.expected, ") Got: ", txt}, "")), nil
}

// EqualIgnoreCase compares the expected value with the incoming request value ignoring the case.
func EqualIgnoreCase(expected string) Matcher {
	return &equalIgnoreCaseMatcher{expected: expected}
}

// EqualIgnoreCasef compares the expected value with the incoming request value ignoring the case.
// This is short-hand to format the expected value.
func EqualIgnoreCasef(format string, a ...any) Matcher {
	return EqualIgnoreCase(fmt.Sprintf(format, a...))
}

// Eqi is an alias to EqualIgnoreCase.
// It compares the expected value with the incoming request value ignoring the case.
func Eqi(expected string) Matcher {
	return &equalIgnoreCaseMatcher{expected: expected}
}
