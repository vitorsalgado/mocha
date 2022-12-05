package matcher

type isMatcher struct {
	matcher Matcher
}

func (m *isMatcher) Name() string {
	return m.matcher.Name()
}

func (m *isMatcher) Match(v any) (*Result, error) {
	return m.matcher.Match(v)
}

func (m *isMatcher) OnMockServed() error {
	return m.matcher.OnMockServed()
}

func Is(matcher Matcher) Matcher {
	return &isMatcher{matcher: matcher}
}
