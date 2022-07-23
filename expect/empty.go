package expect

import "fmt"

// ToBeEmpty returns true if matcher value has zero length.
func ToBeEmpty() Matcher {
	m := Matcher{}
	m.Name = "Empty"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("%v is not empty", v)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		return ToHaveLen(0).Matches(v, args)
	}

	return m
}
