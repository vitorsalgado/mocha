package matcher

type shouldMatcher struct {
	Matcher Matcher
}

func (m *shouldMatcher) Name() string {
	return m.Matcher.Name()
}

func (m *shouldMatcher) Match(v any) (*Result, error) {
	return m.Matcher.Match(v)
}

func (m *shouldMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

func Should(matcher Matcher) Matcher {
	return &shouldMatcher{Matcher: matcher}
}
