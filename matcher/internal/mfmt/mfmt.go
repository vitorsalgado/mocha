package mfmt

import (
	"fmt"
	"net/http"
	"reflect"
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
	return fmt.Sprintf("Received: %s", Stringify(val))
}

type stringer interface {
	String() string
}

func Stringify(v any) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case []byte:
		return string(s)
	case string:
		return s
	case *http.Request:
		return "*http.Request"
	case nil:
		return ""
	default:
		t := reflect.TypeOf(s)
		switch t.Kind() {
		case reflect.Pointer:
			return t.String()
		}
	}

	return fmt.Sprintf("%v", v)
}
