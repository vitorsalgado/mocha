package jsonenc

import (
	"encoding/json"
)

func Encode(matcher string, args ...any) ([]byte, error) {
	switch len(args) {
	case 0:
		return json.Marshal([]any{matcher})
	case 1:
		return json.Marshal([]any{matcher, args[1]})
	}

	d := make([]any, 0, 1+len(args))
	d = append(d, matcher)
	d = append(d, args...)

	return json.Marshal(d)
}
