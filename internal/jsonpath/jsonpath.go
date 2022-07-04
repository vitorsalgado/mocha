// Package jsonpath expose functions to retrieve JSON fields value based on its key path.
package jsonpath

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	fieldRegExp = regexp.MustCompile(`(\w+)\[(\d+)](.*)`)
	idxRegExp   = regexp.MustCompile(`^\[(\d+)](.*)`)

	// ErrFieldNotFound is thrown when a JSON path is invalid.
	ErrFieldNotFound = errors.New("could not find a field using provided json path")
)

// Get returns the field with the given path from the given json data
// Example:
//	data := map[string]any{"address": map[string]any{"street": "somewhere"}}
//	Get("address.street", data)
//	will return "somewhere"
// For arrays, use the notation "field[index]" to get a specific index
func Get(chain string, data any) (any, error) {
	var dataType = reflect.TypeOf(data).Kind()
	var hasBracket = strings.HasPrefix(chain, "[")

	switch dataType {
	case reflect.Map:
		if hasBracket {
			return nil,
				fmt.Errorf("json path starts with an index pattern [n] when the json is actually an object")
		}

		chain := strings.TrimPrefix(chain, ".")

		if matches := fieldRegExp.FindAllStringSubmatch(chain, -1); len(matches) > 0 {
			values := matches[0]
			idx, _ := strconv.Atoi(values[2])

			field := data.(map[string]any)[values[1]].([]any)
			size := len(field)

			if idx > size-1 {
				return nil, ErrFieldNotFound
			}

			entry := field[idx]
			next := values[3]

			if next != "" {
				return Get(next, entry)
			}

			if entry == nil {
				if len(field) <= (idx + 1) {
					return nil, nil
				}

				return nil, ErrFieldNotFound
			}

			return entry, nil
		}

		parts := strings.Split(chain, ".")
		s := len(parts)

		if s == 0 {
			return nil, fmt.Errorf("invalid json path %s", chain)
		}

		if s == 1 {
			val, ok := data.(map[string]any)[parts[0]]
			if !ok {
				return nil, ErrFieldNotFound
			}

			return val, nil
		}

		ch := strings.Join(parts[1:], ".")
		val, ok := data.(map[string]any)[parts[0]]
		if !ok {
			return nil, ErrFieldNotFound
		}

		return Get(ch, val)

	case reflect.Slice:
		if !hasBracket {
			return nil,
				fmt.Errorf("json is an array but the json path does not start with an index pattern []")
		}

		if matches := idxRegExp.FindAllStringSubmatch(chain, -1); len(matches) > 0 {
			values := matches[0]
			idx, _ := strconv.Atoi(values[1])

			d := data.([]any)
			size := len(d)

			if idx > size-1 {
				return nil, ErrFieldNotFound
			}

			next := values[2]
			field := data.([]any)[idx]

			if next == "" {
				if field == nil && len(data.([]any)) < (idx+1) {
					return nil, ErrFieldNotFound
				}

				return field, nil
			}

			return Get(next, field)
		}
	}

	return nil, ErrFieldNotFound
}
