package matcher

import (
	"fmt"
	"regexp"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type regExpMatcher struct {
	expression any
}

func (m *regExpMatcher) Name() string {
	return "RegExp"
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
			fmt.Errorf("matcher does not accept the expression of type %T", v)
	}
}

// Matches passes when the given regular expression matches the incoming request value.
// It accepts a string, regexp.Regexp or *regexp.Regexp.
func Matches(expression any) Matcher {
	return &regExpMatcher{expression: expression}
}
