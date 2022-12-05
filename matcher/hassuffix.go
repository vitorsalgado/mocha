package matcher

import (
	"fmt"
	"strings"
)

type hasSuffixMatcher struct {
	Suffix string
}

func (m *hasSuffixMatcher) Name() string {
	return "HasSuffix"
}

func (m *hasSuffixMatcher) Match(v any) (*Result, error) {
	txt := v.(string)

	return &Result{
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

func (m *hasSuffixMatcher) OnMockServed() error {
	return nil
}

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(suffix string) Matcher {
	return &hasSuffixMatcher{Suffix: suffix}
}
