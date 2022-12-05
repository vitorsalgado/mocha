package matcher

type shouldMatcher struct {
	matcher Matcher
}

func (m *shouldMatcher) Name() string {
	return m.matcher.Name()
}

func (m *shouldMatcher) Match(v any) (*Result, error) {
	return m.matcher.Match(v)
}

func (m *shouldMatcher) OnMockServed() error {
	return m.matcher.OnMockServed()
}

func Should(matcher Matcher) Matcher {
	return &shouldMatcher{matcher: matcher}
}
