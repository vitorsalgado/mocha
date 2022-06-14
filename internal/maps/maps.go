package maps

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var fieldRegexp = regexp.MustCompile("(?P<field>\\w+)\\[(?P<idx>\\d+)]")

func Get[T any](p string, m map[string]any) T {
	var parts = strings.Split(p, ".")
	var ret T

	for i, part := range parts {
		var matches = fieldRegexp.FindAllStringSubmatch(part, -1)
		var field any
		var fieldType reflect.Type

		if matches != nil && len(matches) > 0 {
			values := matches[0]
			keys := fieldRegexp.SubexpNames()
			params := make(map[string]string, len(values))

			for c, key := range values {
				params[keys[c]] = key
			}

			idx, _ := strconv.Atoi(params["idx"])
			f := params["field"]
			field = m[f]
			field = field.([]any)[idx]
			fieldType = reflect.TypeOf(field)
		} else {
			field = m[part]
			fieldType = reflect.TypeOf(field)
		}

		if field == nil {
			var noop T
			return noop
		}

		switch fieldType.Kind() {
		case reflect.Map:
			n := strings.Join(parts[i+1:], ".")
			return Get[T](n, field.(map[string]any))
		default:
			ret = field.(T)
		}
	}

	return ret
}
