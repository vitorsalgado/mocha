package matcher

import (
	"github.com/vitorsalgado/mocha/v3/internal/jsonx"
)

type hasKeyMatcher struct {
	path string
}

func (m *hasKeyMatcher) Name() string {
	return "HasKey"
}

func (m *hasKeyMatcher) Match(v any) (*Result, error) {
	value, err := jsonx.Reach(m.path, v)
	if err != nil {
		return nil, err
	}

	return &Result{
		Pass:    value != nil,
		Message: hint(m.Name(), printExpected(m.path)),
	}, nil
}

func (m *hasKeyMatcher) After() error {
	return nil
}

// HaveKey returns true if the JSON key in the given path is present.
// Example:
//
//	JSON: { "name": "test" }
//	HaveKey("name") will return true
//	HaveKey("address.street") will return false.
func HaveKey(path string) Matcher {
	return &hasKeyMatcher{path: path}
}
