package matcher

import (
	"strconv"
)

type truthyMatcher struct {
}

func (m *truthyMatcher) Name() string {
	return "Truthy"
}

func (m *truthyMatcher) Match(v any) (*Result, error) {
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
		return nil, err
	}

	return &Result{Pass: b}, nil
}

func Truthy() Matcher {
	return &truthyMatcher{}
}
