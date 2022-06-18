package jsonpath

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	fieldRegExp = regexp.MustCompile("(\\w+)\\[(\\d+)](.*)")
	idxRegExp   = regexp.MustCompile("^\\[(\\d+)](.*)")
)

func Get[R any](chain string, data any) (R, error) {
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

		var matches = fieldRegExp.FindAllStringSubmatch(chain, -1)

		if matches != nil && len(matches) > 0 {
			values := matches[0]

			idx, _ := strconv.Atoi(values[2])
			field := data.(map[string]any)[values[1]]
			field = field.([]any)[idx]
			next := values[3]

			if next != "" {
				return Get[R](next, field)
			}

			if field == nil {
				return ret, nil
			}

			return field.(R), nil
		}

		parts := strings.Split(chain, ".")
		s := len(parts)

		if s == 0 {
			return ret, fmt.Errorf("invalid json path %s", chain)
		}

		if s == 1 {
			r := data.(map[string]any)[parts[0]]
			if r == nil {
				return ret, nil
			}

			return data.(map[string]any)[parts[0]].(R), nil
		}

		return Get[R](strings.Join(parts[1:], "."), data.(map[string]any)[parts[0]])

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
				if field == nil {
					return ret, nil
				}

				return field.(R), nil
			}

			return Get[R](next, field)
		}
	}

	return ret, nil
}
