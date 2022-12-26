package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/types"
)

type lowerCaseMatcher struct {
	matcher Matcher
}

func (m *lowerCaseMatcher) Name() string {
	return "ToLower"
}

func (m *lowerCaseMatcher) Match(v any) (*Result, error) {
	txt := v.(string)
	result, err := m.matcher.Match(strings.ToLower(txt))
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		Pass: result.Pass,
		Message: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), printExpected(txt)),
				result.Message(),
			)
		},
	}, nil
}

func (m *lowerCaseMatcher) AfterMockSent() error {
	return m.matcher.AfterMockSent()
}

func (m *lowerCaseMatcher) Raw() types.RawValue {
	return types.RawValue{_mLowerCase, m.matcher.Raw()}
}

// ToLower lower case matcher string argument before submitting it to provided matcher.
func ToLower(matcher Matcher) Matcher {
	return &lowerCaseMatcher{matcher: matcher}
}
