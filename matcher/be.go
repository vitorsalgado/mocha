package matcher

type beMatcher struct {
	matcher Matcher
}

func (m *beMatcher) Name() string {
	return m.matcher.Name()
}

func (m *beMatcher) Match(v any) (*Result, error) {
	return m.matcher.Match(v)
}

func (m *beMatcher) OnMockServed() error {
	return m.matcher.OnMockServed()
}

func Be(matcher Matcher) Matcher {
	return &beMatcher{matcher: matcher}
}
