package expect

type PeekMatcher struct {
	Matcher Matcher
	Action  func(v any) error
}

func (m *PeekMatcher) Name() string {
	return "Peek -> " + m.Matcher.Name()
}

func (m *PeekMatcher) Match(v any) (bool, error) {
	err := m.Action(v)
	if err != nil {
		return false, err
	}

	return m.Matcher.Match(v)
}

func (m *PeekMatcher) DescribeFailure(v any) string {
	return m.Matcher.DescribeFailure(v)
}

func (m *PeekMatcher) OnMockServed() {
	m.Matcher.OnMockServed()
}

// Peek will return the result of the given matcher, after executing the provided function.
// Peek can be used to check the matcher argument.
func Peek(matcher Matcher, action func(v any) error) Matcher {
	return &PeekMatcher{Matcher: matcher, Action: action}
}
