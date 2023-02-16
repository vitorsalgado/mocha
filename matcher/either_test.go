package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEither(t *testing.T) {
	testCases := []struct {
		name     string
		matchers []Matcher
		expected bool
	}{
		{"left true", []Matcher{StrictEqual("test"), Contain("qa")}, true},
		{"right true", []Matcher{StrictEqual("qa"), Contain("tes")}, true},
		{"both true", []Matcher{StrictEqual("test"), Contain("te")}, true},
		{"both false", []Matcher{StrictEqual("qa"), Contain("dev")}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Either(tc.matchers[0], tc.matchers[1]).Match("test")
			require.Nil(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestEitherLeftErr(t *testing.T) {
	result, err := Either(
		Func(func(_ any) (bool, error) {
			return false, fmt.Errorf("fail")
		}),
		Contain("qa")).Match("test")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestEitherRightErr(t *testing.T) {
	result, err := Either(
		Contain("qa"),
		Func(func(_ any) (bool, error) {
			return false, fmt.Errorf("fail")
		})).Match("test")

	require.Error(t, err)
	require.Nil(t, result)
}
