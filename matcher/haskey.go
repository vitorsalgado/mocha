package matcher

import (
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
	if err != nil {
		return mismatch(nil), err
	}

	return Result{
		OK: value != nil,
		DescribeFailure: func() string {
			return hint(m.Name(), printExpected(m.Path))
		},
	}, nil
}

func (m *HasKeyMatcher) OnMockServed() error {
	return nil
}

// HaveKey returns true if the JSON key in the given path is present.
// Example:
//
//	JSON: { "name": "test" }
//	HaveKey("name") will return true
//	HaveKey("address.street") will return false.
func HaveKey(path string) Matcher {
	return &HasKeyMatcher{Path: path}
}
