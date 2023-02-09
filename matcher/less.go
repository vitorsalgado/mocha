package matcher

import "fmt"

type lessMatcher struct {
	expected float64
}

func (m *lessMatcher) Name() string {
	return "Less"
}

func (m *lessMatcher) Match(v any) (*Result, error) {
	vv, err := convToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf("unhandled data type. %w", err)
	}

	if vv < m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: fmt.Sprintf(
		"%s %s %v",
		hint(m.Name(), printExpected(m.expected)),
		_separator,
		printReceived(vv))}, nil
}

func LessThan(expected float64) Matcher {
	return &lessMatcher{expected: expected}
}
