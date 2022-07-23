package expect

import "fmt"

// Repeat returns true if total request hits for current mock is equal or lower total the provided max call times.
// If Repeat is used direct, it must be set using Mock After Expectations.
func Repeat(times int) Matcher {
	count := 0

	m := Matcher{}
	m.Name = "Repeat"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("")
	}
	m.Matches = func(_ any, params Args) (bool, error) {
		count++

		return count <= times, nil
	}

	return m
}
