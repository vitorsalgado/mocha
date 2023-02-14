package matcher

import (
	"errors"
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
		return nil, errors.New("unknown typeof value")
	}

	switch listType.Kind() {
	case reflect.String:
		if pass := strings.Contains(listValue.String(), sub.String()); pass {
			return &Result{Pass: true}, nil
		}

		return &Result{
			Message: stringify(listValue),
			Ext:     []string{stringify(m.expected)},
		}, nil
	case reflect.Map:
		keys := listValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if equalValues(keys[i].Interface(), m.expected) {
				return &Result{Pass: true}, nil
			}
		}

		return &Result{Message: stringify(listValue)}, nil
	}

	for i := 0; i < listValue.Len(); i++ {
		if equalValues(listValue.Index(i).Interface(), sub.Interface()) {
			return &Result{Pass: true}, nil
		}
	}

	return &Result{Message: stringify(listValue)}, nil
}

// Contain returns true when the items value is contained in the matcher argument.
func Contain(expected any) Matcher {
	return &containsMatcher{expected: expected}
}

// Containf returns true when the items value is contained in the matcher argument.
func Containf(format string, a ...any) Matcher {
	return &containsMatcher{expected: fmt.Sprintf(format, a...)}
}
