package matcher

import (
	"fmt"
	"strings"

	"github.com/ryanuber/go-glob"
)

type globMatcher struct {
	pattern string
}

func (m *globMatcher) Match(v any) (Result, error) {
	text, ok := v.(string)
	if !ok {
		return Result{}, fmt.Errorf("glob: it only works with string types. got %T", v)
	}

	if glob.Glob(m.pattern, text) {
		return Result{Pass: true}, nil
	}

	return Result{Message: strings.Join([]string{"Glob(", m.pattern, ") ", text}, "")}, nil
}

func GlobMatch(pattern string) Matcher {
	return &globMatcher{pattern}
}
