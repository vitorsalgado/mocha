package matcher

import (
	"fmt"
	"strings"
)

type LowerCaseMatcher struct {
	Matcher Matcher
}

func (m *LowerCaseMatcher) Name() string {
	return "ToLower"
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

// ToLower lower case matcher string argument before submitting it to provided matcher.
func ToLower(matcher Matcher) Matcher {
	return &LowerCaseMatcher{Matcher: matcher}
}
