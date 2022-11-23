package expect

import (
	"fmt"
)

var _ Matcher = (*BothMatcher)(nil)

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

func (m *BothMatcher) Match(value any) (bool, error) {
	r1, err := m.First.Match(value)
	if err != nil {
		return false, err
	}

	r2, err := m.Second.Match(value)

	return r1 && r2, err
}

func (m *BothMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("one of the matchers \"%s, %s\" dit not match", m.First.Name(), m.Second.Name())
}

func (m *BothMatcher) OnMockServed() {
}

// And sets the second matcher.
func (ba *BothMatcherBuilder) And(second Matcher) Matcher {
	return &BothMatcher{First: ba.First, Second: second}
}
