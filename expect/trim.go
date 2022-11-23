package expect

import "strings"

type TrimMatcher struct {
	Matcher Matcher
}

func (m *TrimMatcher) Name() string {
	return "Trim"
}

func (m *TrimMatcher) Match(v any) (bool, error) {
	return m.Matcher.Match(strings.TrimSpace(v.(string)))
}

func (m *TrimMatcher) DescribeFailure(v any) string {
	return m.Matcher.DescribeFailure(v)
}

func (m *TrimMatcher) OnMockServed() {
	m.Matcher.OnMockServed()
}

// Trim trims' spaces of matcher argument before submitting it to the given matcher.
func Trim(matcher Matcher) Matcher {
	return &TrimMatcher{Matcher: matcher}
}
