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
		return &Result{Pass: false}, err
	}

	r2, err := m.second.Match(v)
	if err != nil {
		return &Result{Pass: false}, err
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

	return &Result{Pass: false, Message: desc, Ext: []string{m.first.Name(), m.second.Name()}}, nil
}

func (m *eitherMatcher) AfterMockServed() error {
	return runAfterMockServed(m.first, m.second)
}

// Either matches true when any of the two given matchers returns true.
func Either(first Matcher, second Matcher) Matcher {
	return &eitherMatcher{first: first, second: second}
}
