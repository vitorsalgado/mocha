package matcher

import (
	"fmt"
)

type NotMatcher struct {
	Matcher Matcher
}

func (m *NotMatcher) Name() string {
	return "Not"
}

func (m *NotMatcher) Match(v any) (Result, error) {
	result, err := m.Matcher.Match(v)
	return Result{
		OK: !result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf(
				"%s ! %s",
				hint(m.Name(), m.Matcher.Name()),
				result.DescribeFailure(),
			)
		},
	}, err
}

func (m *NotMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	return &NotMatcher{Matcher: matcher}
}
