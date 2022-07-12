package is

import (
	"github.com/vitorsalgado/mocha/expect"
)

// AllOf matches when all the given matchers returns true.
// Example:
//	AllOf(EqualTo("test"),ToEqualFold("test"),ToContains("tes"))
func AllOf(list ...expect.Matcher) expect.Matcher {
	return expect.AllOf(list...)
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),ToEqualFold("TEST"),ToContains("tes"))
func AnyOf(list ...expect.Matcher) expect.Matcher {
	return expect.AnyOf(list...)
}

// Both BothAre matches true when both given matchers evaluates to true.
func Both(first expect.Matcher) *expect.BothMatcherBuilder {
	return expect.Both(first)
}

// Either matches true when any of the two given matchers returns true.
func Either(first expect.Matcher) *expect.EitherMatcherBuilder {
	return expect.Either(first)
}

// Empty IsEmpty returns true if matcher value has zero length.
func Empty() expect.Matcher {
	return expect.ToBeEmpty()
}

// EqualTo returns true if matcher value is equal to the given parameter value.
func EqualTo(expected any) expect.Matcher {
	return expect.ToEqual(expected)
}

// EqualFold returns true if expected value is equal to matcher value, ignoring case.
// EqualFold uses strings.EqualFold function.
func EqualFold(expected string) expect.Matcher {
	return expect.ToEqualFold(expected)
}

// Present IsPresent checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func Present() expect.Matcher {
	return expect.ToBePresent()
}

// Not negates the provided matcher.
func Not(matcher expect.Matcher) expect.Matcher {
	return expect.Not(matcher)
}

// XOR is an exclusive or matcher
func XOR(first expect.Matcher, second expect.Matcher) expect.Matcher {
	return expect.XOR(first, second)
}
