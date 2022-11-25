package expect

import "fmt"

type EmptyMatcher struct {
}

func (m *EmptyMatcher) Name() string {
	return "Empty"
}

func (m *EmptyMatcher) Match(v any) (Result, error) {
	result, err := ToHaveLen(0).Match(v)
	if err != nil {
		return Result{}, err
	}

	return Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s %s", hint(m.Name()), _separator, v)
		},
	}, nil
}

func (m *EmptyMatcher) OnMockServed() error {
	return nil
}

// ToBeEmpty returns true if matcher value has zero length.
func ToBeEmpty() Matcher {
	return &EmptyMatcher{}
}
