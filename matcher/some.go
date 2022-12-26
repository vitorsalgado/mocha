package matcher

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/types"
)

type someMatcher struct {
	items []any
}

func (m *someMatcher) Name() string {
	return "Some"
}

func (m *someMatcher) Match(v any) (*Result, error) {
	for _, item := range m.items {
		res, err := Equal(v).Match(item)
		if err != nil {
			return nil, err
		}

		if res.Pass {
			return res, nil
		}
	}

	return &Result{
		Message: func() string {
			return fmt.Sprintf(
				"%s %s value %v is not contained in the %v",
				hint(m.Name(), m.items),
				_separator,
				printReceived(v),
				printExpected(m.items),
			)
		}}, nil
}

func (m *someMatcher) AfterMockSent() error {
	return nil
}

func (m *someMatcher) Raw() types.RawValue {
	return types.RawValue{_mSome, m.items}
}

func Some(items ...any) Matcher {
	return &someMatcher{items: items}
}
