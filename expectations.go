package mocha

import (
	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/expect"
)

type Expectation[V any] struct {
	picker expect.ValueSelector[V]
}

func Expect[V any](picker expect.ValueSelector[V]) Expectation[V] {
	return Expectation[V]{picker: picker}
}

func (e Expectation[V]) ToEqual(v V) core.Expectation[V] {
	return expectation(e.picker, expect.ToEqual(v))
}

func expectation[V any](
	picker expect.ValueSelector[V],
	matcher expect.Matcher[V],
) core.Expectation[V] {
	return core.Expectation[V]{
		Name:          "any",
		ValueSelector: picker,
		Matcher:       matcher,
		Weight:        weightNone,
	}
}
