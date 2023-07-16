package matcher

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type eachMatcher struct {
	matcher Matcher
}

func (m *eachMatcher) Match(v any) (Result, error) {
	var val = reflect.ValueOf(v)
	var valType = reflect.TypeOf(v).Kind()

	switch valType {
	case reflect.Map:
		iterator := val.MapRange()

		for iterator.Next() {
			mv := iterator.Value().Interface()
			res, err := m.matcher.Match(mv)
			if err != nil {
				return Result{}, fmt.Errorf("each: %w", err)
			}

			if !res.Pass {
				return Result{
					Message: strings.
						Join([]string{
							"Each(",
							mfmt.Stringify(iterator.Key().Interface()), ":", mfmt.Stringify(mv),
							") ", res.Message}, "")}, nil
			}
		}

		return Result{Pass: true}, nil

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			entry := val.Index(i).Interface()
			res, err := m.matcher.Match(entry)
			if err != nil {
				return Result{}, fmt.Errorf("each: %w", err)
			}

			if !res.Pass {
				return Result{
					Message: strings.
						Join([]string{
							"Each(", mfmt.Stringify(entry), ") ", res.Message}, "")}, nil
			}
		}

		return Result{Pass: true}, nil
	}

	return Result{}, fmt.Errorf("type %s is not supported. accepted types: map, array", valType.String())
}

func (m *eachMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// Each applies the given matcher on all items of the incoming request value.
// It works with slices and maps.
func Each(matcher Matcher) Matcher {
	return &eachMatcher{matcher: matcher}
}
