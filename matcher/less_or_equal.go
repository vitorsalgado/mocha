package matcher

import "fmt"

type lessOrEqualMatcher struct {
	expected float64
}

func (m *lessOrEqualMatcher) Name() string {
	return "LessOrEqual"
}

func (m *lessOrEqualMatcher) Match(v any) (*Result, error) {
	vv, ok := v.(float64)
	if !ok {
		return nil, fmt.Errorf("matcher Less only works with float64 type")
	}

	if vv <= m.expected {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: fmt.Sprintf(
		"%s %s %v",
		hint(m.Name(), printExpected(m.expected)),
		_separator,
		printReceived(vv))}, nil
}

func (m *lessOrEqualMatcher) After() error {
	return nil
}

func LessOrEqual(expected float64) Matcher {
	return &lessOrEqualMatcher{expected: expected}
}
