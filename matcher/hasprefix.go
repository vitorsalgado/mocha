package matcher

import (
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type hasPrefixMatcher struct {
	prefix string
}

func (m *hasPrefixMatcher) Name() string {
	return "HasPrefix"
}

func (m *hasPrefixMatcher) Match(v any) (*Result, error) {
	txt := v.(string)
	if strings.HasPrefix(txt, m.prefix) {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{m.prefix},
		Message: mfmt.PrintReceived(txt),
	}, nil
}

// HasPrefix returns true if the matcher argument starts with the given prefix.
func HasPrefix(prefix string) Matcher {
	return &hasPrefixMatcher{prefix: prefix}
}
