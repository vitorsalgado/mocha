package core

import (
	"fmt"

	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/internal/misc"
)

// debug wraps an expect.Matcher adding debugging logs.
// The return value will be the same of the provided expect.Matcher.
// debug is used internally by Mocha.
func debug[V any](mk *Mock, matcher expect.Matcher[V]) expect.Matcher[V] {
	m := expect.Matcher[V]{}
	m.Name = "debug"
	m.Matches = func(v V, params expect.Args) (bool, error) {
		result, err := matcher.Matches(v, params)
		desc := fmt.Sprintf("mock: %d", mk.ID)
		if mk.Name != "" {
			desc = desc + " - " + mk.Name
		}

		fmt.Printf(desc + "\n")
		fmt.Printf("matcher: %s\n", colorize.Green(colorize.Bold(matcher.Name)))
		fmt.Printf("received: %s\n", colorize.Gray(misc.ToString(v)))

		if result {
			fmt.Printf("result: %s", colorize.Green("ok"))
		} else {
			fmt.Printf("result: %s", colorize.Red("nok"))
		}

		if err != nil {
			fmt.Printf("\nerror: %s", colorize.Red(err.Error()))
		}

		return result, err
	}

	return m
}
