package matcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type equalJSONMatcher struct {
	expected any
}

func (m *equalJSONMatcher) Match(v any) (Result, error) {
	var expJSON []byte
	var err error

	switch vv := m.expected.(type) {
	case string:
		expJSON = []byte(vv)
	default:
		expJSON, err = json.Marshal(m.expected)
		if err != nil {
			return Result{}, fmt.Errorf("equal_json: %w", err)
		}
	}

	var exp any
	err = json.Unmarshal(expJSON, &exp)
	if err != nil {
		return Result{}, fmt.Errorf("equal_json: %w", err)
	}

	switch vv := v.(type) {
	case string:
		b := new(any)
		err = json.Unmarshal([]byte(vv), b)
		if err != nil {
			return Result{}, fmt.Errorf("equal_json: %w", err)
		}

		if equalValues(*b, exp, false) {
			return Result{Pass: true}, nil
		}
	default:
		if equalValues(v, exp, false) {
			return Result{Pass: true}, nil
		}
	}

	return Result{
		Message: strings.Join([]string{"EqualJSON(", mfmt.Stringify(v), ") Got: ", string(expJSON)}, ""),
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
