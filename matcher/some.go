package matcher

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type someMatcher struct {
	matcher Matcher
}

func (m *someMatcher) Match(v any) (Result, error) {
	if v == nil {
		return mismatch("Some() Expected some values to match but got nil"), nil
	}

	typeOfV := reflect.TypeOf(v)
	kind := typeOfV.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return Result{}, fmt.Errorf("some: matcher only works with slices. got: %v", typeOfV)
	}

	vv := reflect.ValueOf(v)
	if vv.IsZero() {
		return mismatch("Some() Expected some values to match but got zero"), nil
	}

	length := vv.Len()
	messages := make([]string, 0, length)

	for i := 0; i < length; i++ {
		res, err := m.matcher.Match(vv.Index(i).Interface())
		if err != nil {
			return Result{}, fmt.Errorf("some: %w", err)
		}

		if res.Pass {
			return Result{Pass: true}, nil
		}

		messages = append(messages, res.Message)
	}

	return Result{
		Message: fmt.Sprintf(
			"Some()\n%s",
			mfmt.Indent(strings.Join(messages, "\n")),
		),
	}, nil
}

func (m *someMatcher) Describe() any {
	return []any{"some", describe(m.matcher)}
}

// Some will use the given matcher to test whether at least one element in the request value passes.
func Some(matcher Matcher) Matcher {
	return &someMatcher{matcher: matcher}
}
