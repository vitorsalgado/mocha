package matcher

import (
	"fmt"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

type hasKeyMatcher struct {
	path string
	expr jp.Expr
}

func (m *hasKeyMatcher) Match(v any) (Result, error) {
	var results []any

	switch vv := v.(type) {
	case string:
		data, err := oj.ParseString(vv)
		if err != nil {
			return Result{}, fmt.Errorf("has_key: %s: error parsing incoming json: %w", m.path, err)
		}

		results = m.expr.Get(data)
	default:
		results = m.expr.Get(vv)
	}

	size := len(results)
	if size == 0 || (size == 1 && results[0] == nil) {
		return Result{
			Message: strings.Join([]string{"HasKey(", m.path, ") Is not present"}, ""),
		}, nil
	}

	return Result{Pass: true}, nil
}

// HasKey passes if the JSON key in the given path is present.
// Example:
//
//	JSON: { "name": "test" }
//	HasKey("name") will pass
//	HasKey("address.street") will not pass.
func HasKey(path string) Matcher {
	x, err := jp.ParseString(path)
	if err != nil {
		panic(fmt.Errorf("the json path expression %s is invalid: %w", path, err))
	}

	return &hasKeyMatcher{
		path: path,
		expr: x,
	}
}
