package matcher

import (
	"fmt"
	"strings"
)

type hasPrefixMatcher struct {
	prefix string
}

func (m *hasPrefixMatcher) Name() string {
	return "HasPrefix"
}

func (m *hasPrefixMatcher) Match(v any) (*Result, error) {
	txt := v.(string)

	return &Result{
		Pass: strings.HasPrefix(txt, m.prefix),
		Message: func() string {
			return fmt.Sprintf(
				"%s %s %s",
				hint(m.Name(), printExpected(m.prefix)),
				_separator,
				txt,
			)
		},
	}, nil
}

func (m *hasPrefixMatcher) OnMockServed() error {
	return nil
}

func (m *hasPrefixMatcher) Spec() any {
	return []any{_mHasPrefix, m.prefix}
}

// HasPrefix returns true if the matcher argument starts with the given prefix.
func HasPrefix(prefix string) Matcher {
	return &hasPrefixMatcher{prefix: prefix}
}
