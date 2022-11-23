package expect

import (
	"fmt"
	"reflect"
	"strings"
)

var _ Matcher = (*ContainsMatcher)(nil)

type ContainsMatcher struct {
	Expected any
}

func (m *ContainsMatcher) Name() string {
	return "Contains"
}

func (m *ContainsMatcher) Match(list any) (bool, error) {
	listValue := reflect.ValueOf(list)
	sub := reflect.ValueOf(m.Expected)
	listType := reflect.TypeOf(list)
	if listType == nil {
		return false, nil
	}

	kind := listType.Kind()

	if kind == reflect.String {
		return strings.Contains(listValue.String(), sub.String()), nil
	}

	if kind == reflect.Map {
		keys := listValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if reflect.DeepEqual(keys[i].Interface(), m.Expected) {
				return true, nil
			}
		}

		return false, nil
	}

	for i := 0; i < listValue.Len(); i++ {
		if reflect.DeepEqual(listValue.Index(i).Interface(), sub.Interface()) {
			return true, nil
		}
	}

	return false, nil
}

func (m *ContainsMatcher) DescribeFailure(value any) string {
	return fmt.Sprintf("value %v is not contained on %v", m.Expected, value)
}

func (m *ContainsMatcher) OnMockServed() {
}

// ToContain returns true when the expected value is contained in the matcher argument.
func ToContain(expected any) Matcher {
	return &ContainsMatcher{Expected: expected}
}
