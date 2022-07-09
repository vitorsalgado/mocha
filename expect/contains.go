package expect

import "strings"

// ToContain returns true when the expected value is contained in the matcher argument.
func ToContain(expectation string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "ToContains"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.Contains(v, expectation), nil
	}

	return m
}
