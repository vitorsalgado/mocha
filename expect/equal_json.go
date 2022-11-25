package expect

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type EqualJSONMatcher struct {
	Expected any
}

func (m *EqualJSONMatcher) Name() string {
	return "EqualJSON"
}

func (m *EqualJSONMatcher) Match(v any) (Result, error) {
	expectedAsJson, err := json.Marshal(m.Expected)
	if err != nil {
		return mismatch(nil), err
	}

	var exp any
	err = json.Unmarshal(expectedAsJson, &exp)
	if err != nil {
		return mismatch(nil), err
	}

	return Result{
		OK: reflect.DeepEqual(v, exp),
		DescribeFailure: func() string {
			return fmt.Sprintf("%s\nExpected:\n%s\nReceived:\n%s",
				hint(m.Name()),
				printExpected(m.Expected),
				printReceived(v),
			)
		},
	}, nil
}

func (m *EqualJSONMatcher) OnMockServed() error {
	return nil
}

// ToEqualJSON returns true if matcher value is equal to the given parameter value.
func ToEqualJSON(expected any) Matcher {
	return &EqualJSONMatcher{Expected: expected}
}
