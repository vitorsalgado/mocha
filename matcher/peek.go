package matcher

import "fmt"

type peekMatcher struct {
	matcher Matcher
	action  func(v any) error
}

func (m *peekMatcher) Name() string {
	return fmt.Sprintf("Peek(%s)", m.matcher.Name())
}

func (m *peekMatcher) Match(v any) (*Result, error) {
	err := m.action(v)
	if err != nil {
		return nil, err
	}

	return m.matcher.Match(v)
}

func (m *peekMatcher) After() error {
	return m.matcher.After()
}

// Peek will return the result of the given matcher, after executing the provided function.
// Peek can be used to check the matcher argument.
func Peek(matcher Matcher, action func(v any) error) Matcher {
	return &peekMatcher{matcher: matcher, action: action}
}
