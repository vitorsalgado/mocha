package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/types"
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
		Message: func() string {
			return fmt.Sprintf(
				"%s %s %s",
				hint(m.Name(), printExpected(m.suffix)),
				_separator,
				txt,
			)
		},
	}, nil
}

func (m *hasSuffixMatcher) AfterMockSent() error {
	return nil
}

func (m *hasSuffixMatcher) Raw() types.RawValue {
	return types.RawValue{_mHasSuffix, m.suffix}
}

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(suffix string) Matcher {
	return &hasSuffixMatcher{suffix: suffix}
}
