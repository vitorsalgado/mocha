package matcher

import (
	"fmt"
	"reflect"
)

type convertMatcher struct {
	to      reflect.Type
	matcher Matcher
}

func (m *convertMatcher) Name() string {
	return "Convert"
}

func (m *convertMatcher) Match(v any) (*Result, error) {
	vValue := reflect.ValueOf(v)

	if !vValue.CanConvert(m.to) {
		return nil, fmt.Errorf("incoming value %v is convertible to the type %s", v, vValue.Type().Name())
	}

	converted := vValue.Convert(m.to)
	res, err := m.matcher.Match(converted.Interface())
	if err != nil {
		return nil, err
	}

	if res.Pass {
		return res, nil
	}

	return &Result{
		Message: fmt.Sprintf("(%s) %s", stringify(v), res.Message),
		Ext:     []string{stringify(m.to)},
	}, nil
}

func ConvertTo[T any](matcher Matcher) Matcher {
	var t T
	return &convertMatcher{to: reflect.TypeOf(t), matcher: matcher}
}
