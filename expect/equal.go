package expect

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/internal/misc"
)

type EqualMatcher struct {
	Expected any
}

func (m *EqualMatcher) Name() string {
	return "Equal"
}

func (m *EqualMatcher) Match(v any) (bool, error) {
	return reflect.DeepEqual(m.Expected, v), nil
}

func (m *EqualMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("%s\n%s",
		fmt.Sprintf("expected: %v", colorize.Green(misc.Stringify(m.Expected))),
		fmt.Sprintf("got: %s", colorize.Yellow(misc.Stringify(v))),
	)
}

func (m *EqualMatcher) OnMockServed() error {
	return nil
}

// ToEqual returns true if matcher value is equal to the given parameter value.
func ToEqual(expected any) Matcher {
	return &EqualMatcher{Expected: expected}
}
