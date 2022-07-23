package expect

import (
	"fmt"
	"strings"
)

// ToHaveSuffix returns true when matcher argument ends with the given suffix.
func ToHaveSuffix(suffix string) Matcher {
	m := Matcher{}
	m.Name = "HasSuffix"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("value %v, doest not have the suffix %s", v, suffix)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		return strings.HasSuffix(v.(string), suffix), nil
	}

	return m
}
