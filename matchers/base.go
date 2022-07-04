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

func (m Matcher[V]) Describe(describe string) {
	m.Description = describe
}

func (m Matcher[V]) Describef(format string, a ...any) {
	m.Description = fmt.Sprintf(format, a...)
}

func emptyArgs() Args {
	return Args{}
}
