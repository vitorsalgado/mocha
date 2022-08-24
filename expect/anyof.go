package expect

import (
	"fmt"
	"strings"
)

// AnyOf matches when any of the given matchers returns true.
// Example:
//
//	AnyOf(EqualTo("test"),ToEqualFold("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	m := Matcher{}
	m.Name = "AnyOf"
	m.DescribeMismatch = func(p string, v any) string {
		b := make([]string, len(matchers))
		for i, matcher := range matchers {
			b[i] = matcher.Name
		}

		return fmt.Sprintf("none of the given matchers \"%s\" matched.", strings.Join(b, ","))
	}
	m.Matches = func(v any, args Args) (bool, error) {
		for _, matcher := range matchers {
			if result, err := matcher.Matches(v, args); result || err != nil {
				return result, err
			}
		}

		return false, nil
	}

	return m
}
