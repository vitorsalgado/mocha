package expect

import (
	"fmt"
	"strings"
)

var _ Matcher = (*AllOfMatcher)(nil)

type AllOfMatcher struct {
	Matchers []Matcher

	failures []string
}

func (m *AllOfMatcher) Name() string {
	return "AllOf"
}

func (m *AllOfMatcher) Match(v any) (bool, error) {
	for _, matcher := range m.Matchers {
		if result, err := matcher.Match(v); !result || err != nil {
			return result, err
		}
	}

	return true, nil
}

func (m *AllOfMatcher) DescribeFailure(_ any) string {
	b := make([]string, len(m.Matchers))
	for i, matcher := range m.Matchers {
		b[i] = matcher.Name()
	}

	return fmt.Sprintf("one or more of the matchers \"%s\" did not match.", strings.Join(b, ","))
}

func (m *AllOfMatcher) OnMockServed() {
}

// AllOf matches when all the given matchers returns true.
// Example:
//
//	AllOf(EqualTo("test"),ToEqualFold("test"),ToContains("tes"))
func AllOf(matchers ...Matcher) Matcher {
	return &AllOfMatcher{Matchers: matchers}
}
