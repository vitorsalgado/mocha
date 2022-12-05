package matcher

type isMatcher struct {
	Matcher Matcher
}

func (m *isMatcher) Name() string {
	return m.Matcher.Name()
}

func (m *isMatcher) Match(v any) (*Result, error) {
	return m.Matcher.Match(v)
}

func (m *isMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

func Is(matcher Matcher) Matcher {
	return &isMatcher{Matcher: matcher}
}
