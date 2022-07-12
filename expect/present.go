package expect

import "reflect"

// ToBePresent checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func ToBePresent() Matcher {
	m := Matcher{}
	m.Name = "Present"
	m.Matches = func(v any, args Args) (bool, error) {
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

	return m
}
