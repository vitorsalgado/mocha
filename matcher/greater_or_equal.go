package matcher

import (
	"fmt"
)

type greaterOrEqualMatcher struct {
	expected float64
}

func (m *greaterOrEqualMatcher) Name() string {
	return "Less"
}

func (m *greaterOrEqualMatcher) Match(v any) (*Result, error) {
	vv, err := convToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf("unhandled data type. %w", err)
	}

	if vv >= m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{stringify(m.expected)},
		Message: fmt.Sprintf("Received: %v", vv),
	}, nil
}

func GreaterOrEqualThan(expected float64) Matcher {
	return &greaterOrEqualMatcher{expected: expected}
}
