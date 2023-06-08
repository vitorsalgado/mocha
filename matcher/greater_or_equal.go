package matcher

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mconv"
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type greaterOrEqualMatcher struct {
	expected float64
}

func (m *greaterOrEqualMatcher) Name() string {
	return "GreaterOrEqual"
}

func (m *greaterOrEqualMatcher) Match(v any) (*Result, error) {
	vv, err := mconv.ConvToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf(
			"type %T is not supported. the value must be compatible with float64. %w",
			v,
			err,
		)
	}

	if vv >= m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{mfmt.Stringify(m.expected)},
		Message: mfmt.PrintReceived(vv),
	}, nil
}

// GreaterThanOrEqual passes if the incoming request value is greater than or equal to the given value.
func GreaterThanOrEqual(expected float64) Matcher {
	return &greaterOrEqualMatcher{expected: expected}
}

// Gte passes if the incoming request value is greater than or equal to the given value.
func Gte(expected float64) Matcher {
	return &greaterOrEqualMatcher{expected: expected}
}
