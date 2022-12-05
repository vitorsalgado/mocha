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
		return &Result{OK: false}, err
	}

	r2, err := m.second.Match(value)
	if err != nil {
		return &Result{OK: false}, err
	}

	msg := func() string {
		desc := ""

		if !r1.OK {
			desc = r1.DescribeFailure()
		}

		if !r2.OK {
			desc += "\n\n"
			desc += r2.DescribeFailure()
		}

		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), m.first.Name(), m.second.Name()),
			_separator,
			desc)
	}

	return &Result{OK: r1.OK && r2.OK, DescribeFailure: msg}, nil
}

func (m *bothMatcher) OnMockServed() error {
	return multiOnMockServed(m.first, m.second)
}

// Both matches true when both given matchers evaluates to true.
func Both(first Matcher, second Matcher) Matcher {
	m := &bothMatcher{first: first, second: second}
	return m
}
