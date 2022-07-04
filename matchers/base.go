package matchers

import "net/http"

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

	Matcher[V any] func(v V, params Args) (bool, error)
)
