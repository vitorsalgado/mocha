package matcher

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type itemMatcher struct {
	index   int
	matcher Matcher
}

func (m *itemMatcher) Match(v any) (Result, error) {
	if m.index < 0 {
		return Result{}, fmt.Errorf("item: index must greater or equal to 0. got %d", m.index)
	}

	typeOfV := reflect.TypeOf(v)
	kind := typeOfV.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return Result{}, fmt.Errorf("item: %d: it only works with slices. got: %T", m.index, v)
	}

	vv := reflect.ValueOf(v)
	if vv.Len() == 0 {
		return Result{Message: strings.Join([]string{"Item(", strconv.Itoa(m.index), ") Empty array"}, "")}, nil
	}

	res, err := m.matcher.Match(vv.Index(m.index).Interface())
	if err != nil {
		return Result{}, fmt.Errorf("item: %d: %w", m.index, err)
	}

	if res.Pass {
		return Result{Pass: true}, nil
	}

	return Result{
		Message: strings.Join([]string{"Item(", strconv.Itoa(m.index), ") ", res.Message}, ""),
	}, nil
}

func (m *itemMatcher) Describe() any {
	return []any{"item", m.index, describe(m.matcher)}
}

// Item matches a specific array item from the incoming request value.
func Item(index int, matcher Matcher) Matcher {
	return &itemMatcher{index: index, matcher: matcher}
}
