package matcher

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mconv"
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type greaterOrEqualMatcher struct {
	expected float64
}

func (m *greaterOrEqualMatcher) Name() string {
	return "Less"
}

func (m *greaterOrEqualMatcher) Match(v any) (*Result, error) {
	vv, err := mconv.ConvToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf(
			"type %s is not supported. value must be compatible with float64. %w",
			reflect.TypeOf(v),
			err,
		)
	}

	if vv >= m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{mfmt.Stringify(m.expected)},
		Message: fmt.Sprintf("Received: %v", vv),
	}, nil
}

func GreaterOrEqualThan(expected float64) Matcher {
	return &greaterOrEqualMatcher{expected: expected}
}
