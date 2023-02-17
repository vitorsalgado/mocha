package matcher

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type regExpMatcher struct {
	expression any
}

func (m *regExpMatcher) Name() string {
	return "MatchRegExp"
}

func (m *regExpMatcher) Match(v any) (*Result, error) {
	txt := fmt.Sprintf("%v", v)
	msg := mfmt.PrintReceived(txt)
	ext := []string{mfmt.Stringify(m.expression)}

	switch e := m.expression.(type) {
	case string:
		match, err := regexp.Match(e, []byte(txt))
		return &Result{Pass: match, Ext: ext, Message: msg}, err
	case regexp.Regexp:
		return &Result{Pass: e.Match([]byte(txt)), Ext: ext, Message: msg}, nil
	case *regexp.Regexp:
		return &Result{Pass: e.Match([]byte(txt)), Ext: ext, Message: msg}, nil
	default:
		return nil,
			fmt.Errorf("matcher does not accept the expression of type %s", reflect.TypeOf(v))
	}
}

// Matches returns true then the given regular expression matches matcher argument.
// It accepts a string or a regexp.Regexp.
func Matches(expression any) Matcher {
	return &regExpMatcher{expression: expression}
}
