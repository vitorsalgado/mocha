package matcher

import "fmt"

type greaterOrEqualMatcher struct {
	expected float64
}

func (m *greaterOrEqualMatcher) Name() string {
	return "Less"
}

func (m *greaterOrEqualMatcher) Match(v any) (*Result, error) {
	vv, ok := v.(float64)
	if !ok {
		return nil, fmt.Errorf("matcher Less only works with float64 type")
	}

	if vv >= m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: fmt.Sprintf(
		"%s %s %v",
		hint(m.Name(), printExpected(m.expected)),
		_separator,
		printReceived(vv))}, nil
}

func GreaterOrEqualThan(expected float64) Matcher {
	return &greaterOrEqualMatcher{expected: expected}
}
