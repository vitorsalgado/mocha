package matcher

import (
	"errors"
	"strings"
)

func multiOnMockServed(matchers ...Matcher) error {
	var errs []string

	for _, matcher := range matchers {
		err := matcher.AfterMockSent()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}
