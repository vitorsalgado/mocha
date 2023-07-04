package mfmt

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

var _indentRegExp = regexp.MustCompile("(?m)^")

func Indent(str string) string {
	return _indentRegExp.ReplaceAllString(str, " ")
}

func PrintReceived(val any) string {
	return fmt.Sprintf("received: %s", Stringify(val))
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
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(s)
	case *http.Request:
		return "*http.Request"
	case nil:
		return "<nil>"
	default:
		t := reflect.TypeOf(s)
		switch t.Kind() {
		case reflect.Pointer:
			return t.String()
		}
	}

	return fmt.Sprintf("%v", v)
}
