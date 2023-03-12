package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func equalValues(expected any, actual any, lenient bool) bool {
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
	actualValue := reflect.ValueOf(actual)

	if actualKind != expectedKind && lenient {
		pass := laxEq(actualKind, expectedKind, actualValue, expectedValue, lenient)
		if pass {
			return true
		}

		pass = laxEq(expectedKind, actualKind, expectedValue, actualValue, lenient)
		if pass {
			return true
		}

		return false
	}

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

			if !equalValues(b, a, lenient) {
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

			if !equalValues(expectedEntry.Interface(), actualEntry.Interface(), lenient) {
				return false
			}
		}

		return true
	}

	return false
}

func laxEq(
	aKind reflect.Kind,
	bKind reflect.Kind,
	aValue reflect.Value,
	bValue reflect.Value,
	lenient bool,
) bool {
	switch aKind {
	case reflect.String:
		switch bKind {
		case reflect.Float64, reflect.Float32:
			a, err := strconv.ParseFloat(aValue.String(), 64)
			if err != nil {
				return false
			}

			return equalValues(a, bValue.Float(), lenient)

		case reflect.Bool:
			a, err := strconv.ParseBool(aValue.String())
			if err != nil {
				return false
			}

			return equalValues(a, bValue.Bool(), lenient)

		case reflect.Int, reflect.Int32, reflect.Int64:
			a, err := strconv.ParseInt(aValue.String(), 10, 64)
			if err != nil {
				return false
			}

			return equalValues(a, bValue.Int(), lenient)
		}
	}

	return false
}

func runAfterMockServed(matchers ...Matcher) error {
	var errs []string

	for _, matcher := range matchers {
		m, ok := matcher.(OnAfterMockServed)
		if !ok {
			continue
		}

		err := m.AfterMockServed()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func prettierName(m Matcher, r *Result) string {
	var ext []string
	if r != nil {
		ext = r.Ext
	}

	return fmt.Sprintf("%s(%s)", m.Name(), strings.Join(ext, ", "))
}
