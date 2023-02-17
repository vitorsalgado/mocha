package mfmt

import (
	"fmt"
	"regexp"
)

var _indentRegExp = regexp.MustCompile("(?m)^")

func Indent(str string) string {
	times := 1
	rep := ""

	for times > 0 {
		rep += " "
		times--
	}

	return _indentRegExp.ReplaceAllString(str, rep)
}

func PrintReceived(val any) string {
	return fmt.Sprintf("Received: %v", val)
}

type stringer interface {
	String() string
}

func Stringify(v any) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	}

	return fmt.Sprintf("%v", v)
}
