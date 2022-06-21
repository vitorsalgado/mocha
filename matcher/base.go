package matcher

import "net/http"

type (
	Extras interface {
		Get(key string) (any, bool)
	}

	RequestInfo struct {
		Request    *http.Request
		ParsedBody any
	}

	Params struct {
		RequestInfo *RequestInfo
		Extras      Extras
	}

	Matcher[V any] func(v V, params Params) (bool, error)
)
