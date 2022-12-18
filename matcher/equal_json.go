package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"
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
		return mismatch(nil), err
	}

	var exp any
	err = json.Unmarshal(expectedAsJson, &exp)
	if err != nil {
		return mismatch(nil), err
	}

	return &Result{
		Pass: reflect.DeepEqual(v, exp),
		Message: func() string {
			return fmt.Sprintf("%s\nExpected:\n%s\nReceived:\n%s",
				hint(m.Name()),
				printExpected(m.expected),
				printReceived(v),
			)
		},
	}, nil
}

func (m *equalJSONMatcher) OnMockServed() error {
	return nil
}

func (m *equalJSONMatcher) Spec() any {
	return []any{_mEqual, m.expected}
}

// EqualJSON returns true if matcher value is equal to the given parameter value.
func EqualJSON(expected any) Matcher {
	return &equalJSONMatcher{expected: expected}
}
