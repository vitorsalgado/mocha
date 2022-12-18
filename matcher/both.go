package matcher

import (
	"fmt"
)

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
		return &Result{Pass: false}, err
	}

	r2, err := m.second.Match(value)
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

	return &Result{Pass: r1.Pass && r2.Pass, Message: msg}, nil
}

func (m *bothMatcher) OnMockServed() error {
	return multiOnMockServed(m.first, m.second)
}

func (m *bothMatcher) Spec() any {
	return []any{_mBoth, []any{m.first.Spec(), m.second.Spec()}}
}

// Both matches true when both given matchers evaluates to true.
func Both(first Matcher, second Matcher) Matcher {
	m := &bothMatcher{first: first, second: second}
	return m
}
