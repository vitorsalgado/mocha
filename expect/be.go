package expect

import "fmt"

// ToBe returns the result of the provided matcher.
func ToBe(matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "ToBe"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("matcher %s did not match", matcher.Name)
	}
	m.Matches = matcher.Matches

	return m
}
