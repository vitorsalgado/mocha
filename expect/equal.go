package expect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/internal/misc"
)

// ToEqual returns true if matcher value is equal to the given parameter value.
func ToEqual(expected any) Matcher {
	matcher := Matcher{}
	matcher.Name = "Equal"
	matcher.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("%s\n%s",
			fmt.Sprintf("expected: %v", colorize.Green(misc.Stringify(expected))),
			fmt.Sprintf("got: %s", colorize.Yellow(misc.Stringify(v))),
		)
	}
	matcher.Matches = func(v any, args Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}

	return matcher
}

// ToEqualJSON returns true if matcher value is equal to the given parameter value.
func ToEqualJSON(expected any) Matcher {
	matcher := Matcher{}
	matcher.Name = "EqualJSON"
	matcher.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("%s\n%s",
			fmt.Sprintf("expected: %v", colorize.Green(misc.Stringify(expected))),
			fmt.Sprintf("got: %s", colorize.Yellow(misc.Stringify(v))),
		)
	}
	matcher.Matches = func(v any, args Args) (bool, error) {
		expectedAsJson, err := json.Marshal(expected)
		if err != nil {
			return false, err
		}

		var exp any
		err = json.Unmarshal(expectedAsJson, &exp)
		if err != nil {
			return false, err
		}

		return reflect.DeepEqual(v, exp), nil
	}

	return matcher
}

// ToEqualFold returns true if expected value is equal to matcher value, ignoring case.
// ToEqualFold uses strings.EqualFold function.
func ToEqualFold(expected string) Matcher {
	m := Matcher{}
	m.Name = "EqualFold"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("%s\n%s",
			fmt.Sprintf("expected: %v", colorize.Green(expected)),
			fmt.Sprintf("got: %s", colorize.Yellow(v.(string))),
		)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		return strings.EqualFold(expected, v.(string)), nil
	}

	return m
}
