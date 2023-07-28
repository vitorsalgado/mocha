package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type containsMatcher struct {
	expected any
}

func (m *containsMatcher) Match(v any) (Result, error) {
	var eValue = reflect.ValueOf(m.expected)
	var vValue = reflect.ValueOf(v)
	var vType = reflect.TypeOf(v)
	if vType == nil {
		return Result{}, errors.New("contains: unknown typeof value")
	}

	switch vType.Kind() {
	case reflect.String:
		if pass := strings.Contains(vValue.String(), eValue.String()); pass {
			return Result{Pass: true}, nil
		}

		goto ret

	case reflect.Map:
		keys := vValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if equalValues(keys[i].Interface(), m.expected, false) {
				return Result{Pass: true}, nil
			}
		}

		goto ret

	case reflect.Slice, reflect.Array:
		for i := 0; i < vValue.Len(); i++ {
			if equalValues(vValue.Index(i).Interface(), eValue.Interface(), false) {
				return Result{Pass: true}, nil
			}
		}
	}

ret:
	return Result{Message: strings.Join([]string{"Contain(", mfmt.Stringify(m.expected), ") Got: ", mfmt.Stringify(v)}, "")}, nil
}

func (m *containsMatcher) Describe() any {
	return []any{"contains", m.expected}
}

// Contain passes when the expected value is contained in the incoming value from the request.
func Contain(expected any) Matcher {
	return &containsMatcher{expected: expected}
}

// Containf passes when the expected value is contained in the incoming value from the request.
// This is shorthand to format the expected value.
func Containf(format string, a ...any) Matcher {
	return &containsMatcher{expected: fmt.Sprintf(format, a...)}
}
