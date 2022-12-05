package matcher

import (
	"fmt"
	"reflect"
)

type equalMatcher struct {
	Expected any
}

func (m *equalMatcher) Name() string {
	return "Equal"
}

func (m *equalMatcher) Match(v any) (*Result, error) {
	return &Result{
		OK: reflect.DeepEqual(m.Expected, v),
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s %v",
				hint(m.Name(), printExpected(m.Expected)),
				_separator,
				printReceived(v),
			)
		},
	}, nil
}

func (m *equalMatcher) OnMockServed() error {
	return nil
}

// Equal returns true if matcher value is equal to the given parameter value.
func Equal(expected any) Matcher {
	return &equalMatcher{Expected: expected}
}
