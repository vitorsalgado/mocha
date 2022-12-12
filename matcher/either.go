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
		return &Result{OK: false}, err
	}

	r2, err := m.second.Match(v)
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

	return &Result{OK: r1.OK || r2.OK, DescribeFailure: msg}, nil
}

func (m *eitherMatcher) OnMockServed() error {
	return multiOnMockServed(m.first, m.second)
}

func (m *eitherMatcher) Spec() any {
	return []any{_mEither, []any{m.first.Spec(), m.second.Spec()}}
}

// Either matches true when any of the two given matchers returns true.
func Either(first Matcher, second Matcher) Matcher {
	return &eitherMatcher{first: first, second: second}
}
