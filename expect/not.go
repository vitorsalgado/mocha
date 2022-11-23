package expect

import "fmt"

type NotMatcher struct {
	Matcher Matcher
}

func (m *NotMatcher) Name() string {
	return "Not"
}

func (m *NotMatcher) Match(v any) (bool, error) {
	result, err := m.Matcher.Match(v)
	return !result, err
}

func (m *NotMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("matcher %s returned true", m.Matcher.Name())
}

func (m *NotMatcher) OnMockServed() {
	m.Matcher.OnMockServed()
}

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	return &NotMatcher{Matcher: matcher}
}
