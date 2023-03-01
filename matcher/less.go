package matcher

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mconv"
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type lessMatcher struct {
	expected float64
}

func (m *lessMatcher) Name() string {
	return "Less"
}

func (m *lessMatcher) Match(v any) (*Result, error) {
	vv, err := mconv.ConvToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf("type %s is not supported. value must be compatible with float64. %w", reflect.TypeOf(v), err)
	}

	if vv < m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{mfmt.Stringify(m.expected)},
		Message: mfmt.PrintReceived(vv),
	}, nil
}

// LessThan passes if the incoming request value is lower than the given value.
func LessThan(expected float64) Matcher {
	return &lessMatcher{expected: expected}
}
