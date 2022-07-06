package matchers

import (
	"fmt"
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
		Matches     func(v V, args Args) (bool, error)
		Description string
		Name        string
	}
)

func (m *Matcher[V]) Describe(describe string) *Matcher[V] {
	m.Description = describe
	return m
}

func (m *Matcher[V]) Describef(format string, a ...any) *Matcher[V] {
	m.Description = fmt.Sprintf(format, a...)
	return m
}

func (m *Matcher[V]) And(and Matcher[V]) Matcher[V] {
	return AllOf(*m, and)
}

func (m *Matcher[V]) Or(or Matcher[V]) Matcher[V] {
	return AnyOf(*m, or)
}

func (m *Matcher[V]) Xor(and Matcher[V]) Matcher[V] {
	return XOR(*m, and)
}
