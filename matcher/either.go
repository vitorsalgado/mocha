package matcher

import (
	"fmt"
)

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

	msg := func() string {
		desc := ""

		if !r1.Pass {
			desc = r1.Message()
		}

		if !r2.Pass {
			desc += "\n\n"
			desc += r2.Message()
		}

		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), m.first.Name(), m.second.Name()),
			_separator,
			desc)
	}

	return &Result{Pass: r1.Pass || r2.Pass, Message: msg}, nil
}

func (m *eitherMatcher) OnMockServed() error {
	return multiOnMockServed(m.first, m.second)
}

// Either matches true when any of the two given matchers returns true.
func Either(first Matcher, second Matcher) Matcher {
	return &eitherMatcher{first: first, second: second}
}
