package matcher

type eitherMatcher struct {
	first  Matcher
	second Matcher
}

func (m *eitherMatcher) Name() string {
	return "Either"
}

func (m *eitherMatcher) Match(v any) (*Result, error) {
	r1, err := m.first.Match(v)
	if err != nil {
		return nil, err
	}

	r2, err := m.second.Match(v)
	if err != nil {
		return nil, err
	}

	if r1.Pass || r2.Pass {
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
		Message: desc,
		Ext:     []string{prettierName(m.first, r1), prettierName(m.second, r2)},
	}, nil
}

func (m *eitherMatcher) AfterMockServed() error {
	return runAfterMockServed(m.first, m.second)
}

// Either passes when any of the two given matchers pass.
func Either(first Matcher, second Matcher) Matcher {
	return &eitherMatcher{first: first, second: second}
}
