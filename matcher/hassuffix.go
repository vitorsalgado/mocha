package matcher

import (
	"fmt"
	"strings"
)

type HasSuffixMatcher struct {
	Suffix string
}

func (m *HasSuffixMatcher) Name() string {
	return "HasSuffix"
}

func (m *HasSuffixMatcher) Match(v any) (Result, error) {
	txt := v.(string)

	return Result{
		OK: strings.HasSuffix(txt, m.Suffix),
		DescribeFailure: func() string {
			return fmt.Sprintf(
				"%s %s %s",
				hint(m.Name(), printExpected(m.Suffix)),
				_separator,
				txt,
			)
		},
	}, nil
}

func (m *HasSuffixMatcher) OnMockServed() error {
	return nil
}

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(suffix string) Matcher {
	return &HasSuffixMatcher{Suffix: suffix}
}
