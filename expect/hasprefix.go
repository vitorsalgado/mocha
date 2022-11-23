package expect

import (
	"fmt"
	"strings"
)

type HasPrefixMatcher struct {
	Prefix string
}

func (m *HasPrefixMatcher) Name() string {
	return "HasPrefix"
}

func (m *HasPrefixMatcher) Match(v any) (bool, error) {
	return strings.HasPrefix(v.(string), m.Prefix), nil
}

func (m *HasPrefixMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("value %v, doest not have the prefix %s", v, m.Prefix)
}

func (m *HasPrefixMatcher) OnMockServed() error {
	return nil
}

// ToHavePrefix returns true if the matcher argument starts with the given prefix.
func ToHavePrefix(prefix string) Matcher {
	return &HasPrefixMatcher{Prefix: prefix}
}
