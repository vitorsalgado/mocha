package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

type containsMatcher struct {
	expected any
}

func (m *containsMatcher) Name() string {
	return "Contain"
}

func (m *containsMatcher) Match(list any) (*Result, error) {
	var listValue = reflect.ValueOf(list)
	var sub = reflect.ValueOf(m.expected)
	var listType = reflect.TypeOf(list)
	if listType == nil {
		return nil, fmt.Errorf("unknown typeof value")
	}

	var describeFailure = func() string {
		return fmt.Sprintf(
			"%s %s %v",
			hint(m.Name(), printExpected(m.expected)),
			_separator,
			printReceived(listValue),
		)
	}

	switch listType.Kind() {
	case reflect.String:
		return &Result{
			Pass:    strings.Contains(listValue.String(), sub.String()),
			Message: describeFailure,
		}, nil
	case reflect.Map:
		keys := listValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if reflect.DeepEqual(keys[i].Interface(), m.expected) {
				return &Result{
					Pass:    true,
					Message: describeFailure,
				}, nil
			}
		}

		return &Result{Message: describeFailure}, nil
	}

	for i := 0; i < listValue.Len(); i++ {
		if reflect.DeepEqual(listValue.Index(i).Interface(), sub.Interface()) {
			return &Result{
				Pass:    true,
				Message: describeFailure,
			}, nil
		}
	}

	return &Result{Message: describeFailure}, nil
}

func (m *containsMatcher) OnMockServed() error {
	return nil
}

// Contain returns true when the items value is contained in the matcher argument.
func Contain(expected any) Matcher {
	return &containsMatcher{expected: expected}
}

// Containf returns true when the items value is contained in the matcher argument.
func Containf(format string, a ...any) Matcher {
	return &containsMatcher{expected: fmt.Sprintf(format, a...)}
}
