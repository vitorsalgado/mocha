package expect

import "fmt"

type XORMatcher struct {
	First  Matcher
	Second Matcher
}

func (m *XORMatcher) Name() string {
	return "XOR"
}

func (m *XORMatcher) Match(v any) (bool, error) {
	a, err := m.First.Match(v)
	if err != nil {
		return false, err
	}

	b, err := m.Second.Match(v)
	if err != nil {
		return false, err
	}

	return a != b, nil
}

func (m *XORMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("matchers \"%s, %s\" did not meet xor condition", m.First.Name(), m.Second.Name())
}

func (m *XORMatcher) OnMockServed() error {
	return nil
}

// XOR is an exclusive or matcher
func XOR(first Matcher, second Matcher) Matcher {
	return &XORMatcher{First: first, Second: second}
}
