package matcher

import (
	"errors"
	"reflect"
	"strings"
)

func equalValues(expected any, actual any) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}

	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		if reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual) {
			return true
		}
	}

	eType := reflect.TypeOf(expected)
	if eType == nil {
		return false
	}

	actualKind := actualType.Kind()
	expectedKind := eType.Kind()

	if actualKind != expectedKind {
		return false
	}

	actualValue := reflect.ValueOf(actual)

	switch actualKind {
	case reflect.Array, reflect.Slice:
		aLen := actualValue.Len()
		bLen := expectedValue.Len()

		if aLen != bLen {
			return false
		}

		for i := 0; i < aLen; i++ {
			a := actualValue.Index(i).Interface()
			b := expectedValue.Index(i).Interface()

			if !equalValues(b, a) {
				return false
			}
		}

		return true

	case reflect.Map:
		actualKeys := actualValue.MapKeys()
		expectedKeys := expectedValue.MapKeys()

		if len(expectedKeys) != len(actualKeys) {
			return false
		}

		expectedKeyType := expectedValue.Type().Key()
		actualKeyType := actualValue.Type().Key()

		if !expectedKeyType.ConvertibleTo(actualKeyType) {
			return false
		}

		for _, expectedKey := range expectedKeys {
			actualKey := expectedKey.Convert(actualKeyType)
			expectedEntry := expectedValue.MapIndex(expectedKey)
			actualEntry := actualValue.MapIndex(actualKey)

			if !actualEntry.IsValid() {
				return false
			}

			if !equalValues(expectedEntry.Interface(), actualEntry.Interface()) {
				return false
			}
		}

		return true
	}

	return false
}

func runAfter(matchers ...Matcher) error {
	var errs []string

	for _, matcher := range matchers {
		m, ok := matcher.(After)
		if !ok {
			continue
		}

		err := m.After()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}
