package matcher

import (
	"fmt"
)

// BothMatcherBuilder is a builder for Both matcher.
// Use .Both() function to create a new Both matcher.
type BothMatcherBuilder struct {
	First Matcher
}

// Both matches true when both given matchers evaluates to true.
func Both(first Matcher) *BothMatcherBuilder {
	m := &BothMatcherBuilder{First: first}
	return m
}

type BothMatcher struct {
	First  Matcher
	Second Matcher
}

func (m *BothMatcher) Name() string {
	return "Both"
}

func (m *BothMatcher) Match(value any) (Result, error) {
	r1, err := m.First.Match(value)
	if err != nil {
		return Result{OK: false}, err
	}

	r2, err := m.Second.Match(value)
	if err != nil {
		return Result{OK: false}, err
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
			hint(m.Name(), m.First.Name(), m.Second.Name()),
			_separator,
			desc)
	}

	return Result{OK: r1.OK && r2.OK, DescribeFailure: msg}, nil
}

func (m *BothMatcher) OnMockServed() error {
	return multiOnMockServed(m.First, m.Second)
}

// And sets the second matcher.
func (ba *BothMatcherBuilder) And(second Matcher) Matcher {
	return &BothMatcher{First: ba.First, Second: second}
}
