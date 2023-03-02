package matcher

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mconv"
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type lessOrEqualMatcher struct {
	expected float64
}

func (m *lessOrEqualMatcher) Name() string {
	return "LessOrEqual"
}

func (m *lessOrEqualMatcher) Match(v any) (*Result, error) {
	vv, err := mconv.ConvToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf(
			"type %s is not supported. value must be compatible with float64. %w",
			reflect.TypeOf(v),
			err,
		)
	}

	if vv <= m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{
			Ext:     []string{mfmt.Stringify(m.expected)},
			Message: mfmt.PrintReceived(vv)},
		nil
}

// LessThanOrEqual passes if the incoming request value is lower than or equal to the given value.
func LessThanOrEqual(expected float64) Matcher {
	return &lessOrEqualMatcher{expected: expected}
}

// Lte passes if the incoming request value is lower than or equal to the given value.
func Lte(expected float64) Matcher {
	return &lessOrEqualMatcher{expected: expected}
}
