package matchers

import (
	"github.com/vitorsalgado/mocha/base"
	"reflect"
)

func EqualTo[E comparable](expected E) base.Matcher[E] {
	return func(v E, ctx base.MatcherContext) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}
