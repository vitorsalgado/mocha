package matcher

import (
	"fmt"
)

type notMatcher struct {
	matcher Matcher
}

func (m *notMatcher) Name() string {
	return "Not"
}

func (m *notMatcher) Match(v any) (*Result, error) {
	result, err := m.matcher.Match(v)
	if err != nil {
		return nil, err
	}

	return &Result{
		Pass: !result.Pass,
		Message: func() string {
			return fmt.Sprintf(
				"%s ! %s",
				hint(m.Name(), m.matcher.Name()),
				result.Message(),
			)
		},
	}, nil
}

func (m *notMatcher) OnMockServed() error {
	return m.matcher.OnMockServed()
}

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	return &notMatcher{matcher: matcher}
}
