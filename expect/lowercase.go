package expect

import (
	"fmt"
	"strings"
)

type LowerCaseMatcher struct {
	Matcher Matcher
}

func (m *LowerCaseMatcher) Name() string {
	return "LowerCase"
}

func (m *LowerCaseMatcher) Match(v any) (Result, error) {
	txt := v.(string)
	result, err := m.Matcher.Match(strings.ToLower(txt))
	if err != nil {
		return Result{}, err
	}

	return Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), printExpected(txt)),
				result.DescribeFailure(),
			)
		},
	}, nil
}

func (m *LowerCaseMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// LowerCase lower case matcher string argument before submitting it to provided matcher.
func LowerCase(matcher Matcher) Matcher {
	return &LowerCaseMatcher{Matcher: matcher}
}
