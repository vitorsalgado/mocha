package expect

import (
	"fmt"
)

// EitherMatcherBuilder is a builder for Either matcher.
// Prefer to use the Either() function.
type EitherMatcherBuilder struct {
	First Matcher
}

// Either matches true when any of the two given matchers returns true.
func Either(first Matcher) *EitherMatcherBuilder {
	return &EitherMatcherBuilder{first}
}

// Or sets the second matcher
func (e *EitherMatcherBuilder) Or(second Matcher) Matcher {
	return &EitherMatcher{First: e.First, Second: second}
}

type EitherMatcher struct {
	First  Matcher
	Second Matcher
}

func (m *EitherMatcher) Name() string {
	return "Either"
}

func (m *EitherMatcher) Match(v any) (Result, error) {
	r1, err := m.First.Match(v)
	if err != nil {
		return Result{OK: false}, err
	}

	r2, err := m.Second.Match(v)
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

	return Result{OK: r1.OK || r2.OK, DescribeFailure: msg}, nil
}

func (m *EitherMatcher) OnMockServed() error {
	return multiOnMockServed(m.First, m.Second)
}
