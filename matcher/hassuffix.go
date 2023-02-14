package matcher

import (
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
	if strings.HasSuffix(txt, m.suffix) {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{m.suffix},
		Message: printReceived(txt),
	}, nil
}

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(suffix string) Matcher {
	return &hasSuffixMatcher{suffix: suffix}
}
