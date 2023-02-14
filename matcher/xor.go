package matcher

type xorMatcher struct {
	first  Matcher
	second Matcher
}

func (m *xorMatcher) Name() string {
	return "XOR"
}

func (m *xorMatcher) Match(v any) (*Result, error) {
	a, err := m.first.Match(v)
	if err != nil {
		return nil, err
	}

	b, err := m.second.Match(v)
	if err != nil {
		return nil, err
	}

	if a.Pass != b.Pass {
		return &Result{Pass: true}, nil
	}

	desc := ""

	if !a.Pass {
		desc = a.Message
	}

	if !b.Pass {
		desc += "\n\n"
		desc += b.Message
	}

	return &Result{
		Ext:     []string{m.first.Name(), m.second.Name()},
		Message: desc,
	}, nil
}

func (m *xorMatcher) AfterMockServed() error {
	return runAfterMockServed(m.first, m.second)
}

// XOR is an exclusive or matcher
func XOR(first Matcher, second Matcher) Matcher {
	return &xorMatcher{first: first, second: second}
}
