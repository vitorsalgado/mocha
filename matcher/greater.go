package matcher

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mconv"
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type greaterMatcher struct {
	expected float64
}

func (m *greaterMatcher) Name() string {
	return "Greater"
}

func (m *greaterMatcher) Match(v any) (*Result, error) {
	vv, err := mconv.ConvToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf(
			"type %T is not supported. the value must be compatible with float64. %w", v, err)
	}

	if vv > m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{mfmt.Stringify(m.expected)},
		Message: mfmt.PrintReceived(vv),
	}, nil
}

// GreaterThan passes if the incoming request value is greater than the given value.
func GreaterThan(expected float64) Matcher {
	return &greaterMatcher{expected: expected}
}

// Gt passes if the incoming request value is greater than the given value.
func Gt(expected float64) Matcher {
	return &greaterMatcher{expected: expected}
}
