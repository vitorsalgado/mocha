package matcher

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type isContainedInMatcher struct {
	items any
}

func (m *isContainedInMatcher) Name() string {
	return "IsIn"
}

func (m *isContainedInMatcher) Match(v any) (*Result, error) {
	typeOfItems := reflect.TypeOf(m.items)
	kind := typeOfItems.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("matcher only works with slices. got: %v", typeOfItems)
	}

	valueOf := reflect.ValueOf(m.items)

	for i := 0; i < valueOf.Len(); i++ {
		if equalValues(v, valueOf.Index(i).Interface(), false) {
			return &Result{Pass: true}, nil
		}
	}

	return &Result{
		Ext: []string{mfmt.Stringify(m.items)},
		Message: fmt.Sprintf(
			"value %v is not contained in the %v",
			v,
			m.items),
	}, nil
}

// IsIn checks if the incoming request value is in the given items.
// Parameter items must be a slice.
func IsIn(items any) Matcher {
	return &isContainedInMatcher{items: items}
}
