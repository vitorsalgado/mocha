package matcher

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/types"
)

type xorMatcher struct {
	first  Matcher
	second Matcher
}

func (m *xorMatcher) Name() string {
	return "XOR"
}

func (m *xorMatcher) Match(v any) (*Result, error) {
	a, err := m.first.Match(v)
	if err != nil {
		return &Result{}, err
	}

	b, err := m.second.Match(v)
	if err != nil {
		return &Result{}, err
	}

	msg := func() string {
		desc := ""

		if !a.Pass {
			desc = a.Message()
		}

		if !b.Pass {
			desc += "\n\n"
			desc += b.Message()
		}

		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), m.first.Name(), m.second.Name()),
			_separator,
			desc)
	}

	return &Result{
		Pass:    a.Pass != b.Pass,
		Message: msg,
	}, nil
}

func (m *xorMatcher) AfterMockSent() error {
	return nil
}

func (m *xorMatcher) Raw() types.RawValue {
	return types.RawValue{_mXOR, []any{m.first.Raw(), m.second.Raw()}}
}

// XOR is an exclusive or matcher
func XOR(first Matcher, second Matcher) Matcher {
	return &xorMatcher{first: first, second: second}
}
