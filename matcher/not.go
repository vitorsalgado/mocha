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

	if !result.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Message: fmt.Sprintf("!(%s)", result.Message),
		Ext:     []string{prettierName(m.matcher, result)},
	}, nil
}

func (m *notMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	return &notMatcher{matcher: matcher}
}
