package matcher

import (
	"errors"

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
		if errors.Is(err, jsonx.ErrKeyNotFound) {
			return &Result{Pass: false}, nil
		}

		return nil, err
	}

	if value != nil {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: m.path}, nil
}

// HasKey passes if the JSON key in the given path is present.
// Example:
//
//	JSON: { "name": "test" }
//	HasKey("name") will pass
//	HasKey("address.street") will not pass.
func HasKey(path string) Matcher {
	return &hasKeyMatcher{path: path}
}
