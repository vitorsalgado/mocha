package matcher

type nilMatcher struct {
}

func (m *nilMatcher) Name() string {
	return "Nil"
}

func (m *nilMatcher) Match(v any) (*Result, error) {
	if v == nil {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: printReceived(v)}, nil
}

func Nil() Matcher {
	return &nilMatcher{}
}
