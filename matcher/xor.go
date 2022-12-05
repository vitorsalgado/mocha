package matcher

import "fmt"

type xorMatcher struct {
	First  Matcher
	Second Matcher
}

func (m *xorMatcher) Name() string {
	return "XOR"
}

func (m *xorMatcher) Match(v any) (*Result, error) {
	a, err := m.First.Match(v)
	if err != nil {
		return &Result{}, err
	}

	b, err := m.Second.Match(v)
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
			hint(m.Name(), m.First.Name(), m.Second.Name()),
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

// XOR is an exclusive or matcher
func XOR(first Matcher, second Matcher) Matcher {
	return &xorMatcher{First: first, Second: second}
}
