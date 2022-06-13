package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func Equal[T comparable](t *testing.T, expected, actual T, msg ...any) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		Fail(t, fmt.Sprintf("expected: %v. got: %v", expected, actual), msg...)
	}
}

func True(t *testing.T, value bool, msg ...any) {
	t.Helper()

	if !value {
		Fail(t, fmt.Sprintf("expected value to be true. got %t", value), msg...)
	}
}

func False(t *testing.T, value bool, msg ...any) {
	t.Helper()

	if value {
		Fail(t, fmt.Sprintf("expected value to be false. got %t", value), msg...)
	}
}

func Nil(t *testing.T, value any, msg ...any) {
	t.Helper()

	if value != nil {
		Fail(t, fmt.Sprintf("expected value to be nil. got %v", value), msg...)
	}
}

func NotNil(t *testing.T, value any, msg ...any) {
	t.Helper()

	if value == nil {
		Fail(t, fmt.Sprintf("expected value to be not nil. got %v", value), msg...)
	}
}

func Fail(t *testing.T, message string, msgAndArgs ...any) {
	t.Helper()

	extra := msgFromArgs(msgAndArgs...)

	if len(extra) > 0 {
		message = fmt.Sprintf("%s\n%v", message, extra)
	}

	t.Errorf(message)
	t.FailNow()
}

func msgFromArgs(extras ...any) string {
	size := len(extras)

	if size == 0 || extras == nil {
		return ""
	}

	if size == 1 {
		msg := extras[0]
		if str, ok := msg.(string); ok {
			return str
		}

		return fmt.Sprintf("%+v", msg)
	}

	if size > 1 {
		return fmt.Sprintf(extras[0].(string), extras[1:]...)
	}

	return ""
}
