package mock

import (
	"fmt"

	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/to"
)

// Debug wraps a to.Matcher adding debugging logs.
// The return value will be the same of the provided to.Matcher.
// Debug is used internally by Mocha.
func Debug[V any](name string, mk Mock, matcher to.Matcher[V]) to.Matcher[V] {
	m := to.Matcher[V]{}
	m.Name = "Debug"
	m.Matches = func(v V, params to.Args) (bool, error) {
		result, err := matcher.Matches(v, params)

		fmt.Printf("\"%s\" received %v\n", name, colorize.Gray(fmt.Sprintf("%v", v)))

		if err != nil {
			fmt.Print(colorize.RedBright(
				fmt.Sprintf("%s - an error ocurred. reason: %s\n", colorize.Bold(fmt.Sprintf("\"%s\"", name)), err.Error())))
			fmt.Printf("mock %d %s\n", mk.ID, mk.Name)
		} else if !result {
			fmt.Print(colorize.RedBright(
				fmt.Sprintf("%s - did not match\n", colorize.Bold(fmt.Sprintf("\"%s\"", name)))))
			fmt.Printf("mock %d %s\n", mk.ID, mk.Name)
		}

		return result, err
	}

	return m
}
