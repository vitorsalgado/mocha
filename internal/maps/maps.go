package maps

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	fieldRegExp = regexp.MustCompile("(?P<field>\\w+)\\[(?P<idx>\\d+)](.*)")
	idxRegExp   = regexp.MustCompile("^\\[(?P<index>\\d+)](.*)")
)

func GetNew[R any](chain string, data any) (R, error) {
	var dataType = reflect.TypeOf(data).Kind()
	var hasBracket = strings.HasPrefix(chain, "[")
	var ret R

	switch dataType {
	case reflect.Map:
		if hasBracket {
			return ret,
				fmt.Errorf("json path chain starts with an index pattern [] when json is actually an object")
		}

		if strings.HasPrefix(chain, ".") {
			chain = chain[1:]
		}

		parts := strings.Split(chain, ".")
		s := len(parts)

		if s == 0 {
			return ret, fmt.Errorf("invalid json path %s", chain)
		}

		part := parts[0]

		var matches = fieldRegExp.FindAllStringSubmatch(part, -1)
		var field any
		var fieldType reflect.Type

		if matches != nil && len(matches) > 0 {
			values := matches[0]

			idx, _ := strconv.Atoi(values[2])
			field = data.(map[string]any)[values[1]]
			field = field.([]any)[idx]
			next := values[3]

			if next != "" {
				return GetNew[R](next, field)
			}
		} else {
			field = data.(map[string]any)[part]
		}

		fieldType = reflect.TypeOf(field)

		if field == nil {
			var noop R
			return noop, nil
		}

		if fieldType.Kind() == reflect.Map {
			return GetNew[R](strings.Join(parts[1:], "."), field)
		}

		ret = field.(R)

	case reflect.Slice:
		if !hasBracket {
			return *new(R),
				fmt.Errorf("json is an array but the json path chain does not start with an index pattern []")
		}

		matches := idxRegExp.FindAllStringSubmatch(chain, -1)

		if matches != nil && len(matches) > 0 {
			values := matches[0]
			idx, _ := strconv.Atoi(values[1])
			next := values[2]
			field := data.([]any)[idx]

			if next == "" {
				return field.(R), nil
			}

			return GetNew[R](next, field)
		}
	}

	return ret, nil
}
