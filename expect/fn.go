package expect

import "fmt"

// Func creates an anonymous Matcher using the given function.
func Func(fn func(v any, a Args) (bool, error)) Matcher {
	m := Matcher{}
	m.Name = "Func"
	m.Matches = fn
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("custom matcher function did not match. value: %v", v)
	}

	return m
}
