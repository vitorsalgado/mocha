package matcher

type ShouldMatcher struct {
	Matcher Matcher
}

func (m *ShouldMatcher) Name() string {
	return m.Matcher.Name()
}

func (m *ShouldMatcher) Match(v any) (Result, error) {
	return m.Matcher.Match(v)
}

func (m *ShouldMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

func Should(matcher Matcher) Matcher {
	return &ShouldMatcher{Matcher: matcher}
}
