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

func (m *containsMatcher) Name() string {
	return "Contains"
}

func (m *containsMatcher) Match(v any) (*Result, error) {
	var eValue = reflect.ValueOf(m.expected)
	var vValue = reflect.ValueOf(v)
	var vType = reflect.TypeOf(v)
	if vType == nil {
		return nil, errors.New("unknown typeof value")
	}

	switch vType.Kind() {
	case reflect.String:
		if pass := strings.Contains(vValue.String(), eValue.String()); pass {
			return &Result{Pass: true}, nil
		}

		return &Result{
			Message: mfmt.Stringify(vValue),
			Ext:     []string{mfmt.Stringify(m.expected)},
		}, nil

	case reflect.Map:
		keys := vValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if equalValues(keys[i].Interface(), m.expected) {
				return &Result{Pass: true}, nil
			}
		}

		return &Result{Message: mfmt.Stringify(vValue)}, nil

	case reflect.Slice, reflect.Array:
		for i := 0; i < vValue.Len(); i++ {
			if equalValues(vValue.Index(i).Interface(), eValue.Interface()) {
				return &Result{Pass: true}, nil
			}
		}
	}

	return &Result{Message: mfmt.Stringify(vValue)}, nil
}

// Contain passes when the expected value is contained in the incoming value from the request.
func Contain(expected any) Matcher {
	return &containsMatcher{expected: expected}
}

// Containf passes when the expected value is contained in the incoming value from the request.
// This is short-hand to format the expected value.
func Containf(format string, a ...any) Matcher {
	return &containsMatcher{expected: fmt.Sprintf(format, a...)}
}
