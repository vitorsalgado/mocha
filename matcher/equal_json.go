package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type equalJSONMatcher struct {
	expected any
}

func (m *equalJSONMatcher) Name() string {
	return "EqualJSON"
}

func (m *equalJSONMatcher) Match(v any) (*Result, error) {
	var expJSON []byte
	var err error

	switch vv := m.expected.(type) {
	case string:
		expJSON = []byte(vv)
	default:
		expJSON, err = json.Marshal(m.expected)
		if err != nil {
			return nil, err
		}
	}

	var exp any
	err = json.Unmarshal(expJSON, &exp)
	if err != nil {
		return nil, err
	}

	switch vv := v.(type) {
	case string:
		b := new(any)
		err = json.Unmarshal([]byte(vv), b)
		if err != nil {
			return nil, err
		}

		if equalValues(*b, exp, false) {
			return &Result{Pass: true}, err
		}
	default:
		if equalValues(v, exp, false) {
			return &Result{Pass: true}, err
		}
	}

	return &Result{
		Ext:     []string{mfmt.Stringify(m.expected)},
		Message: fmt.Sprintf("expected: %s. received: %s", mfmt.Stringify(m.expected), mfmt.Stringify(v)),
	}, nil
}

// EqualJSON matches JSON values.
func EqualJSON(expected any) Matcher {
	return &equalJSONMatcher{expected: expected}
}

// Eqj matches JSON values.
func Eqj(expected any) Matcher {
	return &equalJSONMatcher{expected: expected}
}
