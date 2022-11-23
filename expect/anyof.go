package expect

import (
	"fmt"
	"strings"
)

var _ Matcher = (*AnyOfMatcher)(nil)

type AnyOfMatcher struct {
	Matchers []Matcher
}

func (m *AnyOfMatcher) Name() string {
	return "AnyOf"
}

func (m *AnyOfMatcher) Match(v any) (bool, error) {
	for _, matcher := range m.Matchers {
		if result, err := matcher.Match(v); result || err != nil {
			return result, err
		}
	}

	return false, nil
}

func (m *AnyOfMatcher) DescribeFailure(_ any) string {
	b := make([]string, len(m.Matchers))
	for i, matcher := range m.Matchers {
		b[i] = matcher.Name()
	}

	return fmt.Sprintf("none of the given matchers \"%s\" matched.", strings.Join(b, ","))
}

func (m *AnyOfMatcher) OnMockServed() error { return nil }

// AnyOf matches when any of the given matchers returns true.
// Example:
//
//	AnyOf(EqualTo("test"),ToEqualFold("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	return &AnyOfMatcher{Matchers: matchers}
}
