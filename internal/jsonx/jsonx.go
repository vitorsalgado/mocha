// Package jsonx implements utilities to work with JSON.
package jsonx

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
)

var ErrKeyNotFound = errors.New("field is not present")

// Reach returns the field with the given path from the given json data
// Example:
//
//	data := map[string]any{"address": map[string]any{"street": "somewhere"}}
//	Reach("address.street", data)
//	will return "somewhere"
//
// For arrays, use the notation "field[index]" to get a specific index
func Reach(path string, data any) (any, error) {
	var dataType = reflect.TypeOf(data).Kind()
	var hasBracket = strings.HasPrefix(path, "[")

	switch dataType {
	case reflect.Map:
		if hasBracket {
			return nil,
				fmt.Errorf("JSON path starts with an index pattern [n] when the JSON it is actually an object")
		}

		chain := strings.TrimPrefix(path, ".")

		if matches := fieldRegExp.FindAllStringSubmatch(chain, -1); len(matches) > 0 {
			values := matches[0]
			idx, _ := strconv.Atoi(values[2])

			field := data.(map[string]any)[values[1]].([]any)
			size := len(field)

			if idx > size-1 {
				return nil, ErrKeyNotFound
			}

			entry := field[idx]
			next := values[3]

			if next != "" {
				return Reach(next, entry)
			}

			if entry == nil {
				if len(field) <= (idx + 1) {
					return nil, nil
				}

				return nil, ErrKeyNotFound
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
				return nil, ErrKeyNotFound
			}

			return val, nil
		}

		ch := strings.Join(parts[1:], ".")
		val, ok := data.(map[string]any)[parts[0]]
		if !ok {
			return nil, ErrKeyNotFound
		}

		return Reach(ch, val)

	case reflect.Slice, reflect.Array:
		if !hasBracket {
			return nil,
				fmt.Errorf("json is an array but the json path does not start with an index pattern []")
		}

		if matches := idxRegExp.FindAllStringSubmatch(path, -1); len(matches) > 0 {
			values := matches[0]
			idx, _ := strconv.Atoi(values[1])

			d := data.([]any)
			size := len(d)

			if idx > size-1 {
				return nil, ErrKeyNotFound
			}

			next := values[2]
			field := data.([]any)[idx]

			if next == "" {
				if field == nil && len(data.([]any)) < (idx+1) {
					return nil, ErrKeyNotFound
				}

				return field, nil
			}

			return Reach(next, field)
		}
	}

	return nil, ErrKeyNotFound
}
