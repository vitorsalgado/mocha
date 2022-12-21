package matcher

type anythingMatcher struct {
}

func (m *anythingMatcher) Name() string {
	return "Anything"
}

func (m *anythingMatcher) Match(v any) (*Result, error) {
	return &Result{Pass: true}, nil
}

func (m *anythingMatcher) OnMockServed() error {
	return nil
}

func (m *anythingMatcher) Spec() any {
	return nil
}

func Anything() Matcher {
	return &anythingMatcher{}
}
