package is

import (
	"github.com/vitorsalgado/mocha/matchers"
)

// AllOf matches when all the given matchers returns true.
// Example:
//	AllOf(EqualTo("test"),EqualFold("test"),Contains("tes"))
func AllOf[V any](list ...matchers.Matcher[V]) matchers.Matcher[V] {
	return matchers.AllOf(list...)
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),EqualFold("TEST"),Contains("tes"))
func AnyOf[V any](list ...matchers.Matcher[V]) matchers.Matcher[V] {
	return matchers.AnyOf(list...)
}

// BothAre matches true when both given matchers evaluates to true.
func Both[V any](first matchers.Matcher[V]) *matchers.BothMatcherBuilder[V] {
	return matchers.Both(first)
}

// Either matches true when any of the two given matchers returns true.
func Either[V any](first matchers.Matcher[V]) *matchers.EitherMatcherBuilder[V] {
	return matchers.Either(first)
}

// IsEmpty returns true if matcher value has zero length.
func Empty[V any](_ ...V) matchers.Matcher[V] {
	return matchers.IsEmpty[V]()
}

// EqualTo returns true if matcher value is equal to the given parameter value.
func EqualTo[V any](expected V) matchers.Matcher[V] {
	return matchers.EqualTo(expected)
}

// EqualFold returns true if expected value is equal to matcher value, ignoring case.
// EqualFold uses strings.EqualFold function.
func EqualFold(expected string) matchers.Matcher[string] {
	return matchers.EqualFold(expected)
}

// IsPresent checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func Present[V any]() matchers.Matcher[V] {
	return matchers.IsPresent[V]()
}

// Not negates the provided matcher.
func Not[V any](matcher matchers.Matcher[V]) matchers.Matcher[V] {
	return matchers.Not(matcher)
}

// XOR is a exclusive or matcher
func XOR[V any](first matchers.Matcher[V], second matchers.Matcher[V]) matchers.Matcher[V] {
	return matchers.XOR(first, second)
}
