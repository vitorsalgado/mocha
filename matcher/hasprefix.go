package matcher

import (
	"fmt"
	"strings"
)

type hasPrefixMatcher struct {
	Prefix string
}

func (m *hasPrefixMatcher) Name() string {
	return "HasPrefix"
}

func (m *hasPrefixMatcher) Match(v any) (*Result, error) {
	txt := v.(string)

	return &Result{
		OK: strings.HasPrefix(txt, m.Prefix),
		DescribeFailure: func() string {
			return fmt.Sprintf(
				"%s %s %s",
				hint(m.Name(), printExpected(m.Prefix)),
				_separator,
				txt,
			)
		},
	}, nil
}

func (m *hasPrefixMatcher) OnMockServed() error {
	return nil
}

// HasPrefix returns true if the matcher argument starts with the given prefix.
func HasPrefix(prefix string) Matcher {
	return &hasPrefixMatcher{Prefix: prefix}
}
