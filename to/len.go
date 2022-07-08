package to

import "reflect"

// HaveLen returns true when matcher argument length is equal to the expected value.
func HaveLen[V any](length int) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Len"
	m.Matches = func(v V, args Args) (bool, error) {
		value := reflect.ValueOf(v)
		return value.Len() == length, nil
	}

	return m
}
