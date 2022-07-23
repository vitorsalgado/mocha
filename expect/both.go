package expect

import (
	"fmt"
)

// BothMatcherBuilder is a builder for Both matcher.
// Use .Both() function to create a new Both matcher.
type BothMatcherBuilder struct {
	first Matcher
}

// Both matches true when both given matchers evaluates to true.
func Both(first Matcher) *BothMatcherBuilder {
	m := &BothMatcherBuilder{first: first}

	return m
}

// And sets the second matcher.
func (ba *BothMatcherBuilder) And(second Matcher) Matcher {
	m := Matcher{}
	m.Name = "Both"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("one of the matchers \"%s, %s\" dit not match", ba.first.Name, second.Name)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		r1, err := ba.first.Matches(v, args)
		if err != nil {
			return false, err
		}

		r2, err := second.Matches(v, args)

		return r1 && r2, err
	}

	return m
}
