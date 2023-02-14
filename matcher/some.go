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
		Ext: []string{stringify(m.items)},
		Message: fmt.Sprintf(
			"Value %v is not contained in the %v",
			v,
			m.items),
	}, nil
}

func Some(items ...any) Matcher {
	return &someMatcher{items: items}
}
