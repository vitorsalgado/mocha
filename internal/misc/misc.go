package misc

import (
	"fmt"
	"reflect"
)

// Stringify returns a string representation of the given parameter, if possible.
func Stringify(v any) string {
	switch e := v.(type) {
	case string:
		return e
	case float64, bool:
		return fmt.Sprintf("%v", e)
	default:
		format := "<value omitted: type=%s>"
		str := "not_defined"

		if e == nil {
			return fmt.Sprintf(format, str)
		}

		nm := reflect.TypeOf(v).Name()
		if nm != "" {
			str = nm
		}

		return fmt.Sprintf(format, str)
	}
}
