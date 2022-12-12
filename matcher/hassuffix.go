package matcher

import (
	"fmt"
	"strings"
)

type hasSuffixMatcher struct {
	suffix string
}

func (m *hasSuffixMatcher) Name() string {
	return "HasSuffix"
}

func (m *hasSuffixMatcher) Match(v any) (*Result, error) {
	txt := v.(string)

	return &Result{
		OK: strings.HasSuffix(txt, m.suffix),
		DescribeFailure: func() string {
			return fmt.Sprintf(
				"%s %s %s",
				hint(m.Name(), printExpected(m.suffix)),
				_separator,
				txt,
			)
		},
	}, nil
}

func (m *hasSuffixMatcher) OnMockServed() error {
	return nil
}

func (m *hasSuffixMatcher) Spec() any {
	return []any{_mHasSuffix, m.suffix}
}

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(suffix string) Matcher {
	return &hasSuffixMatcher{suffix: suffix}
}
