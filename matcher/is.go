package matcher

type IsMatcher struct {
	Matcher Matcher
}

func (m *IsMatcher) Name() string {
	return m.Matcher.Name()
}

func (m *IsMatcher) Match(v any) (Result, error) {
	return m.Matcher.Match(v)
}

func (m *IsMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

func Is(matcher Matcher) Matcher {
	return &IsMatcher{Matcher: matcher}
}
