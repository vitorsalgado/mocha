package is

import (
	"github.com/vitorsalgado/mocha/to"
)

// AllOf matches when all the given matchers returns true.
// Example:
//	AllOf(EqualTo("test"),EqualFold("test"),Contains("tes"))
func AllOf[V any](list ...to.Matcher[V]) to.Matcher[V] {
	return to.BeAllOf(list...)
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),EqualFold("TEST"),Contains("tes"))
func AnyOf[V any](list ...to.Matcher[V]) to.Matcher[V] {
	return to.BeAnyOf(list...)
}

// BothAre matches true when both given matchers evaluates to true.
func Both[V any](first to.Matcher[V]) *to.BothMatcherBuilder[V] {
	return to.Both(first)
}

// Either matches true when any of the two given matchers returns true.
func Either[V any](first to.Matcher[V]) *to.EitherMatcherBuilder[V] {
	return to.Either(first)
}

// IsEmpty returns true if matcher value has zero length.
func Empty[V any](_ ...V) to.Matcher[V] {
	return to.BeEmpty[V]()
}

// EqualTo returns true if matcher value is equal to the given parameter value.
func EqualTo[V any](expected V) to.Matcher[V] {
	return to.Equal(expected)
}

// EqualFold returns true if expected value is equal to matcher value, ignoring case.
// EqualFold uses strings.EqualFold function.
func EqualFold(expected string) to.Matcher[string] {
	return to.EqualFold(expected)
}

// IsPresent checks if matcher argument contains a value that is not nil or the zero value for the argument type.
func Present[V any]() to.Matcher[V] {
	return to.BePresent[V]()
}

// Not negates the provided matcher.
func Not[V any](matcher to.Matcher[V]) to.Matcher[V] {
	return to.Not(matcher)
}

// XOR is a exclusive or matcher
func XOR[V any](first to.Matcher[V], second to.Matcher[V]) to.Matcher[V] {
	return to.XOR(first, second)
}
