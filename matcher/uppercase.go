package matcher

import (
	"fmt"
	"strings"
)

type UpperCaseMatcher struct {
	Matcher Matcher
}

func (m *UpperCaseMatcher) Name() string {
	return m.Matcher.Name()
}

func (m *UpperCaseMatcher) Match(v any) (Result, error) {
	txt := v.(string)
	result, err := m.Matcher.Match(strings.ToUpper(txt))
	if err != nil {
		return Result{}, err
	}

	return Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), m.Matcher.Name()),
				result.DescribeFailure(),
			)
		},
	}, nil
}

func (m *UpperCaseMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// ToUpper upper case matcher string argument before submitting it to provided matcher.
func ToUpper(matcher Matcher) Matcher {
	return &UpperCaseMatcher{Matcher: matcher}
}
