package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBothMatcher(t *testing.T) {
	testCases := []struct {
		name     string
		matchers []Matcher
		expected bool
	}{
		{"left true", []Matcher{StrictEqual("test"), Contain("qa")}, false},
		{"right true", []Matcher{StrictEqual("qa"), Contain("tes")}, false},
		{"both true", []Matcher{StrictEqual("test"), Contain("te")}, true},
		{"both false", []Matcher{StrictEqual("qa"), Contain("dev")}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Both(tc.matchers[0], tc.matchers[1]).Match("test")
			require.Nil(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestBothMatcherErr(t *testing.T) {
	result, err := Both(
		Func(func(_ any) (bool, error) {
			return false, fmt.Errorf("fail")
		}),
		Contain("qa")).Match("test")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestBothMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Both(Eq("yes"), Eq("no")).Name())
}
