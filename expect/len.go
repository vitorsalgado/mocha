package expect

import (
	"fmt"
	"reflect"
)

type LenMatcher struct {
	Length int
}

func (m *LenMatcher) Name() string {
	return "Len"
}

func (m *LenMatcher) Match(v any) (bool, error) {
	value := reflect.ValueOf(v)
	return value.Len() == m.Length, nil
}

func (m *LenMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("value does not have the expected length of %d", m.Length)
}

func (m *LenMatcher) OnMockServed() error {
	return nil
}

// ToHaveLen returns true when matcher argument length is equal to the expected value.
func ToHaveLen(length int) Matcher {
	return &LenMatcher{Length: length}
}
