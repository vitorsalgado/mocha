package expect

import (
	"fmt"
	"strings"
)

// ToHavePrefix returns true if the matcher argument starts with the given prefix.
func ToHavePrefix(prefix string) Matcher {
	m := Matcher{}
	m.Name = "HasPrefix"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("value %v, doest not have the prefix %s", v, prefix)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		return strings.HasPrefix(v.(string), prefix), nil
	}

	return m
}
