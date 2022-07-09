package expect

import "reflect"

// ToBePresent checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func ToBePresent[V any]() Matcher[V] {
	m := Matcher[V]{}
	m.Name = "IsPresent"
	m.Matches = func(v V, args Args) (bool, error) {
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
			return !val.IsZero(), nil
		case reflect.Pointer:
			return !val.IsNil(), nil
		}

		return true, nil
	}

	return m
}
