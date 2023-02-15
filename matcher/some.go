package matcher

import (
	"errors"
	"fmt"
	"reflect"
)

type someMatcher struct {
	items any
}

func (m *someMatcher) Name() string {
	return "Some"
}

func (m *someMatcher) Match(v any) (*Result, error) {
	kind := reflect.TypeOf(m.items).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, errors.New("matcher only works with arrays/slices")
	}

	valueOf := reflect.ValueOf(m.items)

	for i := 0; i < valueOf.Len(); i++ {
		if equalValues(v, valueOf.Index(i).Interface()) {
			return &Result{Pass: true}, nil
		}
	}

	return &Result{
		Ext: []string{stringify(m.items)},
		Message: fmt.Sprintf(
			"Value %v is not contained in the %v",
			v,
			m.items),
	}, nil
}

func Some(items any) Matcher {
	return &someMatcher{items: items}
}
