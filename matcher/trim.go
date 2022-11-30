package matcher

import (
	"fmt"
	"strings"
)

type TrimMatcher struct {
	Matcher Matcher
}

func (m *TrimMatcher) Name() string {
	return "Trim"
}

func (m *TrimMatcher) Match(v any) (Result, error) {
	txt := v.(string)
	result, err := m.Matcher.Match(strings.TrimSpace(txt))
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

func (m *TrimMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// Trim trims' spaces of matcher argument before submitting it to the given matcher.
func Trim(matcher Matcher) Matcher {
	return &TrimMatcher{Matcher: matcher}
}
