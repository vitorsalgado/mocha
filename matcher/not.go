package matcher

import (
	"fmt"
)

type notMatcher struct {
	Matcher Matcher
}

func (m *notMatcher) Name() string {
	return "Not"
}

func (m *notMatcher) Match(v any) (*Result, error) {
	result, err := m.Matcher.Match(v)
	return &Result{
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

func (m *notMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	return &notMatcher{Matcher: matcher}
}
