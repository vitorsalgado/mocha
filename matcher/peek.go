package matcher

import "fmt"

type PeekMatcher struct {
	Matcher Matcher
	Action  func(v any) error
}

func (m *PeekMatcher) Name() string {
	return fmt.Sprintf("Peek(%s)", m.Matcher.Name())
}

func (m *PeekMatcher) Match(v any) (Result, error) {
	err := m.Action(v)
	if err != nil {
		return mismatch(nil), err
	}

	return m.Matcher.Match(v)
}

func (m *PeekMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// Peek will return the result of the given matcher, after executing the provided function.
// Peek can be used to check the matcher argument.
func Peek(matcher Matcher, action func(v any) error) Matcher {
	return &PeekMatcher{Matcher: matcher, Action: action}
}
