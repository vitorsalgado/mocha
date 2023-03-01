package matcher

type anythingMatcher struct {
}

func (m *anythingMatcher) Name() string {
	return "Anything"
}

func (m *anythingMatcher) Match(_ any) (*Result, error) {
	return &Result{Pass: true}, nil
}

// Anything is an empty matcher that always passes.
func Anything() Matcher {
	return &anythingMatcher{}
}
