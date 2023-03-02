package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type equalIgnoreCaseMatcher struct {
	expected string
}

func (m *equalIgnoreCaseMatcher) Name() string {
	return "EqualIgnoreCase"
}

func (m *equalIgnoreCaseMatcher) Match(v any) (*Result, error) {
	if v == nil {
		v = ""
	}

	if strings.EqualFold(m.expected, v.(string)) {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Message: mfmt.PrintReceived(v),
		Ext:     []string{mfmt.Stringify(m.expected)},
	}, nil
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
