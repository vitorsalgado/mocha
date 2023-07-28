package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mconv"
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type greaterMatcher struct {
	expected float64
}

func (m *greaterMatcher) Match(v any) (Result, error) {
	vv, err := mconv.ConvToFloat64(v)
	if err != nil {
		return Result{}, fmt.Errorf(
			"gt: type %T is not supported. the value must be compatible with float64. %w", v, err)
	}

	if vv > m.expected {
		return Result{Pass: true}, nil
	}

	return Result{
		Message: strings.Join([]string{"Gt(", mfmt.Stringify(m.expected), ") Value is ", mfmt.Stringify(v)}, ""),
	}, nil
}

func (m *greaterMatcher) Describe() any {
	return []any{"gt", m.expected}
}

// GreaterThan passes if the incoming request value is greater than the given value.
func GreaterThan(expected float64) Matcher {
	return &greaterMatcher{expected: expected}
}

// Gt passes if the incoming request value is greater than the given value.
func Gt(expected float64) Matcher {
	return &greaterMatcher{expected: expected}
}
