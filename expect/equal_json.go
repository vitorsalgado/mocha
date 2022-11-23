package expect

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/internal/misc"
)

type EqualJSONMatcher struct {
	Expected any
}

func (m *EqualJSONMatcher) Name() string {
	return "EqualJSON"
}

func (m *EqualJSONMatcher) Match(v any) (bool, error) {
	expectedAsJson, err := json.Marshal(m.Expected)
	if err != nil {
		return false, err
	}

	var exp any
	err = json.Unmarshal(expectedAsJson, &exp)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(v, exp), nil
}

func (m *EqualJSONMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("%s\n%s",
		fmt.Sprintf("expected: %v", colorize.Green(misc.Stringify(m.Expected))),
		fmt.Sprintf("got: %s", colorize.Yellow(misc.Stringify(v))),
	)
}

func (m *EqualJSONMatcher) OnMockServed() error {
	return nil
}

// ToEqualJSON returns true if matcher value is equal to the given parameter value.
func ToEqualJSON(expected any) Matcher {
	return &EqualJSONMatcher{Expected: expected}
}
