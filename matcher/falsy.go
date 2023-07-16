package matcher

import (
	"fmt"
	"strconv"
)

type falsyMatcher struct {
}

func (m *falsyMatcher) Match(v any) (Result, error) {
	var b bool
	var err error

	switch e := v.(type) {
	case bool:
		b = e
	case string:
		b, err = strconv.ParseBool(e)
	case int:
		b, err = strconv.ParseBool(strconv.FormatInt(int64(e), 10))
	}

	if err != nil {
		return Result{}, fmt.Errorf("falsy: error parsing value to bool")
	}

	if b {
		return Result{Message: "Falsy() Expected false but it is actually true"}, nil
	}

	return Result{Pass: true}, nil
}

// Falsy checks if the incoming request value is false.
func Falsy() Matcher {
	return &falsyMatcher{}
}
