package to

import "strings"

// Contains returns true when the expected value is contained in the matcher argument.
func Contains(expectation string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "Contains"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.Contains(v, expectation), nil
	}

	return m
}
