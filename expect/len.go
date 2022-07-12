package expect

import "reflect"

// ToHaveLen returns true when matcher argument length is equal to the expected value.
func ToHaveLen(length int) Matcher {
	m := Matcher{}
	m.Name = "Len"
	m.Matches = func(v any, args Args) (bool, error) {
		value := reflect.ValueOf(v)
		return value.Len() == length, nil
	}

	return m
}
