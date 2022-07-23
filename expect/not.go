package expect

import "fmt"

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "Not"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("matcher %s returned true", matcher.Name)
	}
	m.Matches = func(v any, params Args) (bool, error) {
		result, err := matcher.Matches(v, params)
		return !result, err
	}

	return m
}
