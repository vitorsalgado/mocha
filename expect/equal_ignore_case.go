package expect

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
)

type EqualFoldMatcher struct {
	Expected string
}

func (m *EqualFoldMatcher) Name() string {
	return "EqualFold"
}

func (m *EqualFoldMatcher) Match(v any) (bool, error) {
	return strings.EqualFold(m.Expected, v.(string)), nil
}

func (m *EqualFoldMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("%s\n%s",
		fmt.Sprintf("expected: %v", colorize.Green(m.Expected)),
		fmt.Sprintf("got: %s", colorize.Yellow(v.(string))),
	)
}

func (m *EqualFoldMatcher) OnMockServed() {
}

// ToEqualFold returns true if expected value is equal to matcher value, ignoring case.
// ToEqualFold uses strings.EqualFold function.
func ToEqualFold(expected string) Matcher {
	return &EqualFoldMatcher{Expected: expected}
}
