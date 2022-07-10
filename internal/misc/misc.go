package misc

import (
	"fmt"
	"reflect"
)

func ToString(v any) string {
	switch e := v.(type) {
	case string:
		return e
	case float64, bool:
		return fmt.Sprintf("%v", e)
	default:
		return fmt.Sprintf("<value omitted: type=%s>", reflect.TypeOf(v).Name())
	}
}
