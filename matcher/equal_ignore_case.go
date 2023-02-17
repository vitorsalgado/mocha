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
		Message: fmt.Sprintf("Received: %v", v),
		Ext:     []string{mfmt.Stringify(m.expected)},
	}, nil
}

// EqualIgnoreCase returns true if items value is equal to matcher value, ignoring case.
func EqualIgnoreCase(expected string) Matcher {
	return &equalIgnoreCaseMatcher{expected: expected}
}

// EqualIgnoreCasef returns true if items value is equal to matcher value, ignoring case.
func EqualIgnoreCasef(format string, a ...any) Matcher {
	return EqualIgnoreCase(fmt.Sprintf(format, a...))
}
