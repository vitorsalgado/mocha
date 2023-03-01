package matcher

import (
	"fmt"
	"reflect"
)

type eachMatcher struct {
	matcher Matcher
}

func (m *eachMatcher) Name() string {
	return "Each"
}

func (m *eachMatcher) Match(v any) (*Result, error) {
	var val = reflect.ValueOf(v)
	var valType = reflect.TypeOf(v).Kind()

	switch valType {
	case reflect.Map:
		iterator := val.MapRange()

		for iterator.Next() {
			mv := iterator.Value().Interface()
			res, err := m.matcher.Match(mv)
			if err != nil {
				return nil, err
			}

			if !res.Pass {
				return &Result{
					Message: res.Message,
					Ext: []string{
						fmt.Sprintf("key=%v, value=%v", iterator.Key().Interface(), mv),
						prettierName(m.matcher, res),
					},
				}, nil
			}
		}

		return &Result{Pass: true}, nil

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			entry := val.Index(i).Interface()
			res, err := m.matcher.Match(entry)
			if err != nil {
				return nil, err
			}

			if !res.Pass {
				return &Result{
					Message: res.Message,
					Ext: []string{
						fmt.Sprintf("index=%d, item=%v", i, entry),
						prettierName(m.matcher, res),
					},
				}, nil
			}
		}

		return &Result{Pass: true}, nil
	}

	return nil, fmt.Errorf("type %s is not supported. accepted types: map, array", valType.String())
}

func (m *eachMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// Each applies the given matcher on all items of the incoming request value.
// It works with slices and maps.
func Each(matcher Matcher) Matcher {
	return &eachMatcher{matcher: matcher}
}
