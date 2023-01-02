package matcher

import (
	"encoding/json"
	"fmt"
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

	return &Result{
		Pass: equalValues(v, exp),
		Message: fmt.Sprintf("%s\nExpected:\n%s\nReceived:\n%s",
			hint(m.Name()),
			printExpected(m.expected),
			printReceived(v)),
	}, nil
}

func (m *equalJSONMatcher) After() error {
	return nil
}

// EqualJSON returns true if matcher value is equal to the given parameter value.
func EqualJSON(expected any) Matcher {
	return &equalJSONMatcher{expected: expected}
}
