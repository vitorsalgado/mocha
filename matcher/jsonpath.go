package matcher

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

type jsonPathMatcher struct {
	path    string
	expr    jp.Expr
	matcher Matcher
	name    string
}

func (m *jsonPathMatcher) Name() string {
	return m.name
}

func (m *jsonPathMatcher) Match(v any) (*Result, error) {
	var results []any

	switch vv := v.(type) {
	case string:
		data, err := oj.ParseString(vv)
		if err != nil {
			return nil, fmt.Errorf("error parsing incoming json: %w", err)
		}

		results = m.expr.Get(data)
	default:
		results = m.expr.Get(vv)
	}

	var r *Result
	var err error

	size := len(results)
	if size == 0 {
		r, err = m.matcher.Match(nil)
	} else if size == 1 {
		r, err = m.matcher.Match(results[0])
	} else {
		r, err = m.matcher.Match(results)
	}

	if err != nil {
		return nil, err
	}

	if r.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{m.path, prettierName(m.matcher, r)},
		Message: r.Message,
	}, nil
}

func (m *jsonPathMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// JSONPath applies the provided matcher to the JSON field value in the given path.
// Example:
//
//	JSONPath("address.city", EqualTo("Santiago"))
func JSONPath(path string, matcher Matcher) Matcher {
	x, err := jp.ParseString(path)
	if err != nil {
		panic(fmt.Errorf("the json path expression %s is invalid: %w", path, err))
	}

	return &jsonPathMatcher{
		path:    path,
		expr:    x,
		matcher: matcher,
		name:    "JSONPath",
	}
}

// Field is an alias for JSONPath.
// It applies the provided matcher to the JSON field value in the given path.
// Example:
//
//	Field("address.city", EqualTo("Santiago"))
func Field(path string, matcher Matcher) Matcher {
	x, err := jp.ParseString(path)
	if err != nil {
		panic(fmt.Errorf("the json path expression %s is invalid: %w", path, err))
	}

	return &jsonPathMatcher{
		path:    path,
		expr:    x,
		matcher: matcher,
		name:    "Field",
	}
}
