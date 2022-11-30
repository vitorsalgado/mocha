package matcher

import (
	"fmt"
	"reflect"
)

type EachMatcher struct {
	Matcher Matcher
}

func (m *EachMatcher) Name() string {
	return "Each"
}

func (m *EachMatcher) Match(v any) (Result, error) {
	var val = reflect.ValueOf(v)
	var valType = reflect.TypeOf(v).Kind()

	switch valType {
	case reflect.Map:
		iter := val.MapRange()

		for iter.Next() {
			mv := iter.Value().Interface()
			res, err := m.Matcher.Match(mv)
			if err != nil {
				return mismatch(nil), err
			}

			if !res.OK {
				return Result{
					OK: false,
					DescribeFailure: func() string {
						return fmt.Sprintf("%s %s %s",
							hint(
								m.Name(),
								fmt.Sprintf("key=%v, value=%v", iter.Key().Interface(), mv)),
							_separator,
							res.DescribeFailure())
					},
				}, nil
			}
		}

		return Result{OK: true}, nil

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			entry := val.Index(i).Interface()
			res, err := m.Matcher.Match(entry)
			if err != nil {
				return mismatch(nil), err
			}

			if !res.OK {
				return Result{
					OK: false,
					DescribeFailure: func() string {
						return fmt.Sprintf("%s %s %s", hint(
							m.Name(),
							fmt.Sprintf("index=%d, item=%v", i, entry)),
							_separator,
							res.DescribeFailure())
					},
				}, nil
			}
		}

		return Result{OK: true}, nil
	}

	return Result{
		OK: false,
		DescribeFailure: func() string {
			return hint(m.Name(), printReceived(fmt.Sprintf("type %s is not supported", valType.String())))
		},
	}, nil
}

func (m *EachMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

func Each(matcher Matcher) Matcher {
	return &EachMatcher{Matcher: matcher}
}
