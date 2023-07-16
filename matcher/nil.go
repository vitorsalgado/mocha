package matcher

type nilMatcher struct {
}

func (m *nilMatcher) Match(v any) (Result, error) {
	if v == nil {
		return Result{Pass: true}, nil
	}

	return Result{Message: "Nil() Value is not nil"}, nil
}

// IsNil passes if the incoming request value is nil.
func IsNil() Matcher {
	return &nilMatcher{}
}
