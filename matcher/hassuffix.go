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
		Pass: strings.HasSuffix(txt, m.suffix),
		Message: fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.suffix)),
			_separator,
			txt),
	}, nil
}

func (m *hasSuffixMatcher) After() error {
	return nil
}

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(suffix string) Matcher {
	return &hasSuffixMatcher{suffix: suffix}
}
