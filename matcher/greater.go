package matcher

import "fmt"

type greaterMatcher struct {
	expected float64
}

func (m *greaterMatcher) Name() string {
	return "Less"
}

func (m *greaterMatcher) Match(v any) (*Result, error) {
	vv, err := convToFloat64(v)
	if err != nil {
		return nil, fmt.Errorf("unhandled data type. %w", err)
	}

	if vv > m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: fmt.Sprintf(
		"%s %s %v",
		hint(m.Name(), printExpected(m.expected)),
		_separator,
		printReceived(vv))}, nil
}

func GreaterThan(expected float64) Matcher {
	return &greaterMatcher{expected: expected}
}
