package matcher

import "fmt"

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

		if !a.OK {
			desc = a.DescribeFailure()
		}

		if !b.OK {
			desc += "\n\n"
			desc += b.DescribeFailure()
		}

		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), m.first.Name(), m.second.Name()),
			_separator,
			desc)
	}

	return &Result{
		OK:              a.OK != b.OK,
		DescribeFailure: msg,
	}, nil
}

func (m *xorMatcher) OnMockServed() error {
	return nil
}

func (m *xorMatcher) Spec() any {
	return []any{_mXOR, []any{m.first.Spec(), m.second.Spec()}}
}

// XOR is an exclusive or matcher
func XOR(first Matcher, second Matcher) Matcher {
	return &xorMatcher{first: first, second: second}
}
