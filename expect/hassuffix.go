package expect

import (
	"fmt"
	"strings"
)

type HasSuffixMatcher struct {
	Suffix string
}

func (m *HasSuffixMatcher) Name() string {
	return "HasSuffix"
}

func (m *HasSuffixMatcher) Match(v any) (bool, error) {
	return strings.HasSuffix(v.(string), m.Suffix), nil
}

func (m *HasSuffixMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("value %v, doest not have the suffix %s", v, m.Suffix)
}

func (m *HasSuffixMatcher) OnMockServed() {
}

// ToHaveSuffix returns true when matcher argument ends with the given suffix.
func ToHaveSuffix(suffix string) Matcher {
	return &HasSuffixMatcher{Suffix: suffix}
}
