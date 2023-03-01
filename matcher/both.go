package matcher

type bothMatcher struct {
	first  Matcher
	second Matcher
}

func (m *bothMatcher) Name() string {
	return "Both"
}

func (m *bothMatcher) Match(value any) (*Result, error) {
	r1, err := m.first.Match(value)
	if err != nil {
		return nil, err
	}

	r2, err := m.second.Match(value)
	if err != nil {
		return nil, err
	}

	if r1.Pass && r2.Pass {
		return &Result{Pass: true}, nil
	}

	desc := ""

	if !r1.Pass {
		desc = r1.Message
	}

	if !r2.Pass {
		desc += "\n\n"
		desc += r2.Message
	}

	return &Result{
		Pass:    false,
		Message: desc,
		Ext:     []string{prettierName(m.first, r1), prettierName(m.second, r2)},
	}, nil
}

func (m *bothMatcher) AfterMockServed() error {
	return runAfterMockServed(m.first, m.second)
}

// Both passes when both the given matchers pass.
func Both(first Matcher, second Matcher) Matcher {
	return &bothMatcher{first: first, second: second}
}
