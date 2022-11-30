package matcher

type BeMatcher struct {
	Matcher Matcher
}

func (m *BeMatcher) Name() string {
	return m.Matcher.Name()
}

func (m *BeMatcher) Match(v any) (Result, error) {
	return m.Matcher.Match(v)
}

func (m *BeMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

func Be(matcher Matcher) Matcher {
	return &BeMatcher{Matcher: matcher}
}
