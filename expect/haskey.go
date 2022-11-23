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

func (m *HasKeyMatcher) Match(v any) (bool, error) {
	value, err := jsonx.Reach(m.Path, v)
	if err != nil || value == nil {
		return false, err
	}

	return true, nil
}

func (m *HasKeyMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("json doest not have a key on path: %s", m.Path)
}

func (m *HasKeyMatcher) OnMockServed() {
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
