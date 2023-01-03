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
	message := fmt.Sprintf("%s %v", hint(m.Name()), v)

	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Interface:
		return &Result{Pass: !val.IsZero(), Message: message}, nil
	case reflect.Pointer:
		return &Result{Pass: !val.IsNil(), Message: message}, nil
	}

	return &Result{Pass: true, Message: message}, nil
}

func (m *bePresentMatcher) AfterMockServed() error {
	return nil
}

// Present checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func Present() Matcher {
	return &bePresentMatcher{}
}
