package matcher

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mconv"
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type greaterMatcher struct {
	expected float64
}

func (m *greaterMatcher) Name() string {
	return "Less"
}

func (m *greaterMatcher) Match(v any) (*Result, error) {
	vv, err := mconv.ConvToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf(
			"type %s is not supported. the value must be compatible with float64. %w",
			reflect.TypeOf(v),
			err,
		)
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
