package expect

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/internal/misc"
)

// ToEqual returns true if matcher value is equal to the given parameter value.
func ToEqual(expected any) Matcher {
	matcher := Matcher{}
	matcher.Name = "EqualTo"
	matcher.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("%s\n%s",
			fmt.Sprintf("expected: %v", colorize.Green(misc.ToString(expected))),
			fmt.Sprintf("got: %s", colorize.Red(misc.ToString(v))),
		)
	}
	matcher.Matches = func(v any, args Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}

	return matcher
}
