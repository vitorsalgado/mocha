package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/types"
)

type anyOfMatcher struct {
	matchers []Matcher
}

func (m *anyOfMatcher) Name() string {
	return "AnyOf"
}

func (m *anyOfMatcher) Match(v any) (*Result, error) {
	ok := false
	errs := make([]string, 0)
	failed := make([]string, 0)

	for _, matcher := range m.matchers {
		result, err := matcher.Match(v)
		if err != nil {
			ok = false
			errs = append(errs, err.Error())
			failed = append(failed, result.Message())

			continue
		}

		if !result.Pass {
			ok = false
			failed = append(failed, result.Message())

			continue
		}

		ok = true
		break
	}

	describeFailure := func() string {
		return fmt.Sprintf(
			"%s\n%s",
			hint(m.Name(), fmt.Sprintf("+%d", len(m.matchers))),
			indent(strings.Join(failed, "\n")),
		)
	}

	if len(errs) > 0 {
		return &Result{
			Pass:    false,
			Message: describeFailure,
		}, fmt.Errorf(strings.Join(errs, "\n"))
	}

	return &Result{Pass: ok, Message: describeFailure}, nil
}

func (m *anyOfMatcher) AfterMockSent() error {
	return multiOnMockServed(m.matchers...)
}

func (m *anyOfMatcher) Raw() types.RawValue {
	args := make([]any, len(m.matchers))

	for i, matcher := range m.matchers {
		args[i] = matcher.Raw()
	}

	return types.RawValue{_mAnyOf, args}
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//
//	AnyOf(EqualTo("test"),EqualIgnoreCase("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	return &anyOfMatcher{matchers: matchers}
}
