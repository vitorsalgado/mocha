package expect

import (
	"fmt"
)

var _ Matcher = (*EitherMatcher)(nil)

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

func (m *EitherMatcher) Match(v any) (bool, error) {
	r1, err := m.First.Match(v)
	if err != nil {
		return false, err
	}

	r2, err := m.Second.Match(v)

	return r1 || r2, err
}

func (m *EitherMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("none of the matchers \"%s, %s\" matched.", m.First.Name(), m.Second.Name())
}

func (m *EitherMatcher) OnMockServed() error {
	return nil
}
