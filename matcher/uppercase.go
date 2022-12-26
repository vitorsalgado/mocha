package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/types"
)

type upperCaseMatcher struct {
	matcher Matcher
}

func (m *upperCaseMatcher) Name() string {
	return m.matcher.Name()
}

func (m *upperCaseMatcher) Match(v any) (*Result, error) {
	txt := v.(string)
	result, err := m.matcher.Match(strings.ToUpper(txt))
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		Pass: result.Pass,
		Message: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), m.matcher.Name()),
				result.Message(),
			)
		},
	}, nil
}

func (m *upperCaseMatcher) AfterMockSent() error {
	return m.matcher.AfterMockSent()
}

func (m *upperCaseMatcher) Raw() types.RawValue {
	return types.RawValue{_mUpperCase, m.matcher.Raw()}
}

// ToUpper upper case matcher string argument before submitting it to provided matcher.
func ToUpper(matcher Matcher) Matcher {
	return &upperCaseMatcher{matcher: matcher}
}
