package matcher

import (
	"reflect"
	"strconv"
	"strings"
)

type lenMatcher struct {
	length int
}

func (m *lenMatcher) Match(v any) (Result, error) {
	value := reflect.ValueOf(v)
	actual := value.Len()
	if actual == m.length {
		return Result{Pass: true}, nil
	}

	return Result{
		Message: strings.
			Join([]string{"Len(", strconv.Itoa(m.length), ") Got: ", strconv.Itoa(actual)}, "")}, nil
}

// Len passes when the expected value length is equal to the incoming request value.
func Len(length int) Matcher {
	return &lenMatcher{length: length}
}
