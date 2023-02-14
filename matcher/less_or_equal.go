package matcher

import "fmt"

type lessOrEqualMatcher struct {
	expected float64
}

func (m *lessOrEqualMatcher) Name() string {
	return "LessOrEqual"
}

func (m *lessOrEqualMatcher) Match(v any) (*Result, error) {
	vv, err := convToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf("unhandled data type. %w", err)
	}

	if vv <= m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{
			Ext:     []string{stringify(m.expected)},
			Message: printReceived(vv)},
		nil
}

func LessOrEqual(expected float64) Matcher {
	return &lessOrEqualMatcher{expected: expected}
}
