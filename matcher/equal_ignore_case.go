package matcher

import (
	"fmt"
	"strings"
)

type EqualFoldMatcher struct {
	Expected string
}

func (m *EqualFoldMatcher) Name() string {
	return "EqualFold"
}

func (m *EqualFoldMatcher) Match(v any) (Result, error) {
	if v == nil {
		v = ""
	}

	return Result{
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

func (m *EqualFoldMatcher) OnMockServed() error {
	return nil
}

// EqualIgnoreCase returns true if expected value is equal to matcher value, ignoring case.
// EqualIgnoreCase uses strings.EqualFold function.
func EqualIgnoreCase(expected string) Matcher {
	return &EqualFoldMatcher{Expected: expected}
}
