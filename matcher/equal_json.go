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
	expectedAsJson, err := json.Marshal(m.expected)
	if err != nil {
		return nil, err
	}

	var exp any
	err = json.Unmarshal(expectedAsJson, &exp)
	if err != nil {
		return nil, err
	}

	if equalValues(v, exp) {
		return &Result{Pass: true}, err
	}

	return &Result{
		Ext:     []string{mfmt.Stringify(m.expected)},
		Message: fmt.Sprintf("Received: %v", v),
	}, nil
}

// EqualJSON returns true if matcher value is equal to the given parameter value.
func EqualJSON(expected any) Matcher {
	return &equalJSONMatcher{expected: expected}
}
