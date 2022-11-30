package matcher

import (
	"fmt"
	"reflect"
)

type EqualMatcher struct {
	Expected any
}

func (m *EqualMatcher) Name() string {
	return "Equal"
}

func (m *EqualMatcher) Match(v any) (Result, error) {
	return Result{
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

func (m *EqualMatcher) OnMockServed() error {
	return nil
}

// Equal returns true if matcher value is equal to the given parameter value.
func Equal(expected any) Matcher {
	return &EqualMatcher{Expected: expected}
}
