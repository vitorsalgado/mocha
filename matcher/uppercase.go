package matcher

import (
	"fmt"
	"strings"
)

type upperCaseMatcher struct {
	Matcher Matcher
}

func (m *upperCaseMatcher) Name() string {
	return m.Matcher.Name()
}

func (m *upperCaseMatcher) Match(v any) (*Result, error) {
	txt := v.(string)
	result, err := m.Matcher.Match(strings.ToUpper(txt))
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), m.Matcher.Name()),
				result.DescribeFailure(),
			)
		},
	}, nil
}

func (m *upperCaseMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// ToUpper upper case matcher string argument before submitting it to provided matcher.
func ToUpper(matcher Matcher) Matcher {
	return &upperCaseMatcher{Matcher: matcher}
}
