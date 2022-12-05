package matcher

import (
	"fmt"
	"strings"
)

type equalIgnoreCaseMatcher struct {
	Expected string
}

func (m *equalIgnoreCaseMatcher) Name() string {
	return "EqualFold"
}

func (m *equalIgnoreCaseMatcher) Match(v any) (*Result, error) {
	if v == nil {
		v = ""
	}

	return &Result{
		OK: strings.EqualFold(m.Expected, v.(string)),
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s %s",
				hint(m.Name(), printExpected(m.Expected)),
				_separator,
				printReceived(v),
			)
		},
	}, nil
}

func (m *equalIgnoreCaseMatcher) OnMockServed() error {
	return nil
}

// EqualIgnoreCase returns true if expected value is equal to matcher value, ignoring case.
// EqualIgnoreCase uses strings.EqualFold function.
func EqualIgnoreCase(expected string) Matcher {
	return &equalIgnoreCaseMatcher{Expected: expected}
}
