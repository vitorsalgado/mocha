package matcher

import (
	"fmt"
	"reflect"
)

type bePresentMatcher struct {
}

func (m *bePresentMatcher) Name() string {
	return "Present"
}

func (m *bePresentMatcher) Match(v any) (*Result, error) {
	if v == nil {
		return &Result{}, nil
	}

	val := reflect.ValueOf(v)
	message := func() string {
		return fmt.Sprintf("%s %v", hint(m.Name()), v)
	}

	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Interface:
		return &Result{OK: !val.IsZero(), DescribeFailure: message}, nil
	case reflect.Pointer:
		return &Result{OK: !val.IsNil(), DescribeFailure: message}, nil
	}

	return &Result{OK: true, DescribeFailure: message}, nil
}

func (m *bePresentMatcher) OnMockServed() error {
	return nil
}

// Present checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func Present() Matcher {
	return &bePresentMatcher{}
}
