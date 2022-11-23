package expect

import (
	"fmt"
	"reflect"
)

type BePresentMatcher struct {
}

func (m *BePresentMatcher) Name() string {
	return "Present"
}

func (m *BePresentMatcher) Match(v any) (bool, error) {
	if v == nil {
		return false, nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Interface:
		return !val.IsZero(), nil
	case reflect.Pointer:
		return !val.IsNil(), nil
	}

	return true, nil
}

func (m *BePresentMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("%v is not present", v)
}

func (m *BePresentMatcher) OnMockServed() {
}

// ToBePresent checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func ToBePresent() Matcher {
	return &BePresentMatcher{}
}
