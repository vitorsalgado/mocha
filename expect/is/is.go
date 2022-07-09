package is

import (
	"github.com/vitorsalgado/mocha/expect"
)

// AllOf matches when all the given matchers returns true.
// Example:
//	AllOf(EqualTo("test"),ToEqualFold("test"),ToContains("tes"))
func AllOf[V any](list ...expect.Matcher[V]) expect.Matcher[V] {
	return expect.AllOf(list...)
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),ToEqualFold("TEST"),ToContains("tes"))
func AnyOf[V any](list ...expect.Matcher[V]) expect.Matcher[V] {
	return expect.AnyOf(list...)
}

// Both BothAre matches true when both given matchers evaluates to true.
func Both[V any](first expect.Matcher[V]) *expect.BothMatcherBuilder[V] {
	return expect.Both(first)
}

// Either matches true when any of the two given matchers returns true.
func Either[V any](first expect.Matcher[V]) *expect.EitherMatcherBuilder[V] {
	return expect.Either(first)
}

// Empty IsEmpty returns true if matcher value has zero length.
func Empty[V any](_ ...V) expect.Matcher[V] {
	return expect.ToBeEmpty[V]()
}

// EqualTo returns true if matcher value is equal to the given parameter value.
func EqualTo[V any](expected V) expect.Matcher[V] {
	return expect.ToEqual(expected)
}

// EqualFold returns true if expected value is equal to matcher value, ignoring case.
// EqualFold uses strings.EqualFold function.
func EqualFold(expected string) expect.Matcher[string] {
	return expect.ToEqualFold(expected)
}

// Present IsPresent checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func Present[V any]() expect.Matcher[V] {
	return expect.ToBePresent[V]()
}

// Not negates the provided matcher.
func Not[V any](matcher expect.Matcher[V]) expect.Matcher[V] {
	return expect.Not(matcher)
}

// XOR is an exclusive or matcher
func XOR[V any](first expect.Matcher[V], second expect.Matcher[V]) expect.Matcher[V] {
	return expect.XOR(first, second)
}
