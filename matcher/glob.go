package matcher

import (
	"fmt"

	"github.com/ryanuber/go-glob"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type globMatcher struct {
	pattern string
}

func (m *globMatcher) Name() string {
	return "Glob"
}

func (m *globMatcher) Match(v any) (*Result, error) {
	text, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("glob only works with string types")
	}

	if glob.Glob(m.pattern, text) {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: mfmt.PrintReceived(text), Ext: []string{m.pattern}}, nil
}

func GlobMatch(pattern string) Matcher {
	return &globMatcher{pattern}
}
