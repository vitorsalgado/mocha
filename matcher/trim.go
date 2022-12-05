package matcher

import (
	"fmt"
	"strings"
)

type trimMatcher struct {
	Matcher Matcher
}

func (m *trimMatcher) Name() string {
	return "Trim"
}

func (m *trimMatcher) Match(v any) (*Result, error) {
	txt := v.(string)
	result, err := m.Matcher.Match(strings.TrimSpace(txt))
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), printExpected(txt)),
				result.DescribeFailure(),
			)
		},
	}, nil
}

func (m *trimMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// Trim trims' spaces of matcher argument before submitting it to the given matcher.
func Trim(matcher Matcher) Matcher {
	return &trimMatcher{Matcher: matcher}
}
