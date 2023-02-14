package matcher

import (
	"fmt"
	"regexp"
)

var _indentRegExp = regexp.MustCompile("(?m)^")

func indent(str string) string {
	times := 1
	rep := ""

	for times > 0 {
		rep += " "
		times--
	}

	return _indentRegExp.ReplaceAllString(str, rep)
}

func printReceived(val any) string {
	return fmt.Sprintf("Received: %v", val)
}
