package expect

import (
	"fmt"
	"strings"
)

type HasPrefixMatcher struct {
	Prefix string
}

func (m *HasPrefixMatcher) Name() string {
	return "HasPrefix"
}

func (m *HasPrefixMatcher) Match(v any) (Result, error) {
	txt := v.(string)

	return Result{
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

func (m *HasPrefixMatcher) OnMockServed() error {
	return nil
}

// ToHavePrefix returns true if the matcher argument starts with the given prefix.
func ToHavePrefix(prefix string) Matcher {
	return &HasPrefixMatcher{Prefix: prefix}
}
