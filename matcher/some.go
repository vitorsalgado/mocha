package matcher

import "fmt"

type someMatcher struct {
	items []any
}

func (m *someMatcher) Name() string {
	return "Some"
}

func (m *someMatcher) Match(v any) (*Result, error) {
	for _, item := range m.items {
		if equalValues(v, item) {
			return &Result{Pass: true}, nil
		}
	}

	return &Result{
		Message: fmt.Sprintf(
			"%s %s value %v is not contained in the %v",
			hint(m.Name(), m.items),
			_separator,
			printReceived(v),
			printExpected(m.items)),
	}, nil
}

func (m *someMatcher) After() error {
	return nil
}

func Some(items ...any) Matcher {
	return &someMatcher{items: items}
}
