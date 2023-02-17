package matcher

import (
	"reflect"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
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

	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Interface:
		return &Result{Pass: !val.IsZero(), Message: mfmt.Stringify(v)}, nil
	case reflect.Pointer:
		return &Result{Pass: !val.IsNil(), Message: mfmt.Stringify(v)}, nil
	}

	return &Result{Pass: true}, nil
}

// Present checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func Present() Matcher {
	return &bePresentMatcher{}
}
