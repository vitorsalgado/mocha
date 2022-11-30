package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

type ContainsMatcher struct {
	Expected any
}

func (m *ContainsMatcher) Name() string {
	return "Contains"
}

func (m *ContainsMatcher) Match(list any) (Result, error) {
	var listValue = reflect.ValueOf(list)
	var sub = reflect.ValueOf(m.Expected)
	var listType = reflect.TypeOf(list)
	if listType == nil {
		return mismatch(nil), fmt.Errorf("unknown typeof value")
	}

	var describeFailure = func() string {
		return fmt.Sprintf(
			"%s %s %v",
			hint(m.Name(), printExpected(m.Expected)),
			_separator,
			printReceived(listValue),
		)
	}

	switch listType.Kind() {
	case reflect.String:
		return Result{
			OK:              strings.Contains(listValue.String(), sub.String()),
			DescribeFailure: describeFailure,
		}, nil
	case reflect.Map:
		keys := listValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if reflect.DeepEqual(keys[i].Interface(), m.Expected) {
				return Result{
					OK:              true,
					DescribeFailure: describeFailure,
				}, nil
			}
		}

		return mismatch(describeFailure), nil
	}

	for i := 0; i < listValue.Len(); i++ {
		if reflect.DeepEqual(listValue.Index(i).Interface(), sub.Interface()) {
			return Result{
				OK:              true,
				DescribeFailure: describeFailure,
			}, nil
		}
	}

	return mismatch(describeFailure), nil
}

func (m *ContainsMatcher) OnMockServed() error {
	return nil
}

// Contain returns true when the expected value is contained in the matcher argument.
func Contain(expected any) Matcher {
	return &ContainsMatcher{Expected: expected}
}
