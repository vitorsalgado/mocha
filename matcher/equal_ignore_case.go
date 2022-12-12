package matcher

import (
	"fmt"
	"strings"
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

	return &Result{
		OK: strings.EqualFold(m.expected, v.(string)),
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s %s",
				hint(m.Name(), printExpected(m.expected)),
				_separator,
				printReceived(v),
			)
		},
	}, nil
}

func (m *equalIgnoreCaseMatcher) OnMockServed() error {
	return nil
}

func (m *equalIgnoreCaseMatcher) Spec() any {
	return []any{_mEqualIgnoreCase, m.expected}
}

// EqualIgnoreCase returns true if items value is equal to matcher value, ignoring case.
func EqualIgnoreCase(expected string) Matcher {
	return &equalIgnoreCaseMatcher{expected: expected}
}

// EqualIgnoreCasef returns true if items value is equal to matcher value, ignoring case.
func EqualIgnoreCasef(format string, a ...any) Matcher {
	return EqualIgnoreCase(fmt.Sprintf(format, a...))
}
