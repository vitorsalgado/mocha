package expect

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

	Matcher struct {
		Name             string
		DescribeMismatch func(p string, v any) string
		Matches          func(v any, args Args) (bool, error)
	}

	ValueSelector func(r *RequestInfo) any
)

func (m Matcher) And(and Matcher) Matcher {
	return AllOf(m, and)
}

func (m Matcher) Or(or Matcher) Matcher {
	return AnyOf(m, or)
}

func (m Matcher) Xor(and Matcher) Matcher {
	return XOR(m, and)
}
