package matcher

import (
	"fmt"
	"strconv"
)

type truthyMatcher struct {
}

func (m *truthyMatcher) Match(v any) (Result, error) {
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
		return Result{}, fmt.Errorf("truthy: error parsing value to bool")
	}

	if !b {
		return Result{Message: "Truthy() Expected true but it is actually false"}, nil
	}

	return Result{Pass: true}, nil
}

// Truthy passes if the request value is true.
func Truthy() Matcher {
	return &truthyMatcher{}
}
