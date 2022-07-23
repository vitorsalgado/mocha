package expect

import "fmt"

// XOR is an exclusive or matcher
func XOR(first Matcher, second Matcher) Matcher {
	m := Matcher{}
	m.Name = "Xor"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("matchers \"%s, %s\" did not meet xor condition", first.Name, second.Name)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		a, err := first.Matches(v, args)
		if err != nil {
			return false, err
		}

		b, err := second.Matches(v, args)
		if err != nil {
			return false, err
		}

		return a != b, nil
	}

	return m
}
