package expect

import (
	"fmt"
	"reflect"
)

// ToHaveLen returns true when matcher argument length is equal to the expected value.
func ToHaveLen(length int) Matcher {
	m := Matcher{}
	m.Name = "Len"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("value does not have the expected length of %d", length)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		value := reflect.ValueOf(v)
		return value.Len() == length, nil
	}

	return m
}
