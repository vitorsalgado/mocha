package expect

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/internal/jsonx"
)

type HasKeyMatcher struct {
	Path string
}

func (m *HasKeyMatcher) Name() string {
	return "HasKey"
}

func (m *HasKeyMatcher) Match(v any) (Result, error) {
	value, err := jsonx.Reach(m.Path, v)
	if err != nil || value == nil {
		return mismatch(nil), err
	}

	return Result{
		OK: true,
		DescribeFailure: func() string {
			return fmt.Sprintf(
				"%s",
				hint(m.Name(), printExpected(m.Path)),
			)
		},
	}, nil
}

func (m *HasKeyMatcher) OnMockServed() error {
	return nil
}

// ToHaveKey returns true if the JSON key in the given path is present.
// Example:
//
//	JSON: { "name": "test" }
//	ToHaveKey("name") will return true
//	ToHaveKey("address.street") will return false.
func ToHaveKey(path string) Matcher {
	return &HasKeyMatcher{Path: path}
}
