package mocha

import (
	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/expect"
)

type Expectation struct {
	picker expect.ValueSelector
}

func Expect(picker expect.ValueSelector) Expectation {
	return Expectation{picker: picker}
}

func (e Expectation) ToEqual(v any) core.Expectation {
	return expectation(e.picker, expect.ToEqual(v))
}

func (e Expectation) ToContain(v any) core.Expectation {
	return expectation(e.picker, expect.ToContain(v))
}

func expectation(
	picker expect.ValueSelector,
	matcher expect.Matcher,
) core.Expectation {
	return core.Expectation{
		Name:          "any",
		ValueSelector: picker,
		Matcher:       matcher,
		Weight:        _weightNone,
	}
}
