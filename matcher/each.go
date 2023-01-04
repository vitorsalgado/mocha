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
		iter := val.MapRange()

		for iter.Next() {
			mv := iter.Value().Interface()
			res, err := m.matcher.Match(mv)
			if err != nil {
				return nil, err
			}

			if !res.Pass {
				return &Result{
					Pass: false,
					Message: fmt.Sprintf("%s %s %s",
						hint(
							m.Name(),
							fmt.Sprintf("key=%v, value=%v", iter.Key().Interface(), mv)),
						_separator,
						res.Message),
				}, nil
			}
		}

		return &Result{Pass: true}, nil

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			entry := val.Index(i).Interface()
			res, err := m.matcher.Match(entry)
			if err != nil {
				return nil, err
			}

			if !res.Pass {
				return &Result{
					Pass: false,
					Message: fmt.Sprintf("%s %s %s", hint(
						m.Name(),
						fmt.Sprintf("index=%d, item=%v", i, entry)),
						_separator,
						res.Message),
				}, nil
			}
		}

		return &Result{Pass: true}, nil
	}

	return &Result{
		Pass:    false,
		Message: hint(m.Name(), printReceived(fmt.Sprintf("type %s is not supported", valType.String()))),
	}, nil
}

func (m *eachMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

func Each(matcher Matcher) Matcher {
	return &eachMatcher{matcher: matcher}
}
