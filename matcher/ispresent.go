package matcher

import "reflect"

func IsPresent[V any]() Matcher[V] {
	return func(v V, params Params) (bool, error) {
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
			return !val.IsZero(), nil
		case reflect.Pointer:
			return !val.IsNil(), nil
		}

		return true, nil
	}
}
