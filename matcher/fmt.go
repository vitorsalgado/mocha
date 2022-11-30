package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
)

var _indentRegExp = regexp.MustCompile("(?m)^")

func hint(name string, extra ...any) string {
	if len(extra) == 0 {
		return colorize.Bold(name)
	}

	ext := make([]string, len(extra))

	for i, e := range extra {
		ext[i] = fmt.Sprintf("%v", e)
	}

	return fmt.Sprintf("%s(%s)", colorize.Bold(name), strings.Join(ext, ", "))
}

func indent(str string) string {
	times := 1
	rep := ""

	for times > 0 {
		rep += " "
		times--
	}

	return _indentRegExp.ReplaceAllString(str, rep)
}

func printExpected(val any) string {
	return colorize.Green(fmt.Sprintf("%v", val))
}

func printReceived(val any) string {
	if val == nil {
		return colorize.Yellow("<nil>")
	}

	return colorize.Yellow(fmt.Sprintf("%v", val))
}
