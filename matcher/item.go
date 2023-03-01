package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type itemMatcher struct {
	index   int
	matcher Matcher
}

func (m *itemMatcher) Name() string {
	return "Item"
}

func (m *itemMatcher) Match(v any) (*Result, error) {
	if m.index < 0 {
		return nil, fmt.Errorf("index must greater or equal to 0. got %d", m.index)
	}

	kind := reflect.TypeOf(v).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, errors.New("matcher only works with arrays/slices")
	}

	vv := reflect.ValueOf(v)

	if vv.Len() == 0 {
		return &Result{Message: "array is empty", Ext: []string{strconv.FormatInt(
			int64(m.index),
			10,
		)}}, nil
	}

	res, err := m.matcher.Match(vv.Index(m.index).Interface())
	if err != nil {
		return nil, err
	}

	if res.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Message: res.Message,
		Ext:     []string{strconv.FormatInt(int64(m.index), 10), prettierName(m.matcher, res)},
	}, nil
}

// Item matches a specific array item from the incoming request value.
func Item(index int, matcher Matcher) Matcher {
	return &itemMatcher{index: index, matcher: matcher}
}
