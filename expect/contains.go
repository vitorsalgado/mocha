package expect

import (
	"reflect"
	"strings"
)

// ToContain returns true when the expected value is contained in the matcher argument.
func ToContain[V any](expectation any) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Contain"
	m.Matches = func(list V, args Args) (bool, error) {
		listValue := reflect.ValueOf(list)
		sub := reflect.ValueOf(expectation)
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
				if reflect.DeepEqual(keys[i].Interface(), list) {
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

	return m
}
