package matcher

import "fmt"

type peekMatcher struct {
	matcher Matcher
	action  func(v any) error
}

func (m *peekMatcher) Match(v any) (Result, error) {
	err := m.action(v)
	if err != nil {
		return Result{}, fmt.Errorf("peek: %w", err)
	}

	return m.matcher.Match(v)
}

// Peek will return the expected of the given matcher, after executing the provided function.
// Peek can be used to check the matcher argument.
func Peek(matcher Matcher, action func(v any) error) Matcher {
	return &peekMatcher{matcher: matcher, action: action}
}
