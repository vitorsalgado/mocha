package to

import (
	"net/http"
)

type (
	Params interface {
		Get(key string) (any, bool)
	}

	RequestInfo struct {
		Request    *http.Request
		ParsedBody any
	}

	Args struct {
		RequestInfo *RequestInfo
		Params      Params
	}

	Matcher[V any] struct {
		Name             string
		DescribeMismatch func(plance string, v any) string
		Matches          func(v V, args Args) (bool, error)
	}
)

func (m *Matcher[V]) And(and Matcher[V]) Matcher[V] {
	return BeAllOf(*m, and)
}

func (m *Matcher[V]) Or(or Matcher[V]) Matcher[V] {
	return BeAnyOf(*m, or)
}

func (m *Matcher[V]) Xor(and Matcher[V]) Matcher[V] {
	return XOR(*m, and)
}
