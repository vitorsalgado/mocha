package matcher

import "fmt"

type lessMatcher struct {
	expected float64
}

func (m *lessMatcher) Name() string {
	return "Less"
}

func (m *lessMatcher) Match(v any) (*Result, error) {
	vv, ok := v.(float64)
	if !ok {
		return nil, fmt.Errorf("matcher Less only works with float64 type")
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

func (m *lessMatcher) After() error {
	return nil
}

func LessThan(expected float64) Matcher {
	return &lessMatcher{expected: expected}
}
