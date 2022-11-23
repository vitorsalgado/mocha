package expect

import "fmt"

var _ Matcher = (*EmptyMatcher)(nil)

type EmptyMatcher struct {
}

func (m *EmptyMatcher) Name() string {
	return "Empty"
}

func (m *EmptyMatcher) Match(v any) (bool, error) {
	return ToHaveLen(0).Match(v)
}

func (m *EmptyMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("%v is not empty", v)
}

func (m *EmptyMatcher) OnMockServed() {
}

// ToBeEmpty returns true if matcher value has zero length.
func ToBeEmpty() Matcher {
	return &EmptyMatcher{}
}
