package matcher

type anythingMatcher struct {
}

func (m *anythingMatcher) Name() string {
	return "Anything"
}

func (m *anythingMatcher) Match(_ any) (*Result, error) {
	return &Result{Pass: true}, nil
}

func (m *anythingMatcher) After() error {
	return nil
}

func Anything() Matcher {
	return &anythingMatcher{}
}
