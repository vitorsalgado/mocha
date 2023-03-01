package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type someMatcher struct {
	matcher Matcher
}

func (m *someMatcher) Name() string {
	return "Some"
}

func (m *someMatcher) Match(v any) (*Result, error) {
	kind := reflect.TypeOf(v).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, errors.New("matcher only works with arrays/slices")
	}

	vv := reflect.ValueOf(v)
	messages := make([]string, 0)

	for i := 0; i < vv.Len(); i++ {
		res, err := m.matcher.Match(vv.Index(i).Interface())
		if err != nil {
			return nil, err
		}

		if res.Pass {
			return &Result{Pass: true}, nil
		}

		messages = append(messages, res.Message)
	}

	return &Result{
		Ext: []string{mfmt.Stringify(v), prettierName(m.matcher, nil)},
		Message: fmt.Sprintf(
			"%s\n%s",
			prettierName(m.matcher, nil),
			mfmt.Indent(strings.Join(messages, "\n")),
		),
	}, nil
}

// Some will use the given matcher to test whether at least one element in the request value passes.
func Some(matcher Matcher) Matcher {
	return &someMatcher{matcher: matcher}
}
