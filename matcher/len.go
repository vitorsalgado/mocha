package matcher

import "reflect"

// Len returns true when matcher argument length is equal to the expected value.
func Len[E any](length int) Matcher[E] {
	return func(v E, params Args) (bool, error) {
		value := reflect.ValueOf(v)
		return value.Len() == length, nil
	}
}
