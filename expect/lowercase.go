package expect

import (
	"strings"
)

type LowerCaseMatcher struct {
	Matcher Matcher
}

func (m *LowerCaseMatcher) Name() string {
	return "LowerCase"
}

func (m *LowerCaseMatcher) Match(v any) (bool, error) {
	return m.Matcher.Match(strings.ToLower(v.(string)))
}

func (m *LowerCaseMatcher) DescribeFailure(v any) string {
	return m.Matcher.DescribeFailure(v)
}

func (m *LowerCaseMatcher) OnMockServed() {
	// TODO implement me
	panic("implement me")
}

// LowerCase lower case matcher string argument before submitting it to provided matcher.
func LowerCase(matcher Matcher) Matcher {
	return &LowerCaseMatcher{Matcher: matcher}
}
