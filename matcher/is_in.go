package matcher

import (
	"errors"
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
	kind := reflect.TypeOf(m.items).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, errors.New("matcher only works with arrays/slices")
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
