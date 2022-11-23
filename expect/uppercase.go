package expect

import "strings"

type UpperCaseMatcher struct {
	Matcher Matcher
}

func (m *UpperCaseMatcher) Name() string {
	return m.Name()
}

func (m *UpperCaseMatcher) Match(v any) (bool, error) {
	return m.Matcher.Match(strings.ToUpper(v.(string)))
}

func (m *UpperCaseMatcher) DescribeFailure(v any) string {
	return m.DescribeFailure(v)
}

func (m *UpperCaseMatcher) OnMockServed() {
	m.OnMockServed()
}

// UpperCase upper case matcher string argument before submitting it to provided matcher.
func UpperCase(matcher Matcher) Matcher {
	return &UpperCaseMatcher{Matcher: matcher}
}
