package mocha

import (
	"github.com/vitorsalgado/mocha/expect"
)

func Header(name string) expect.ValueSelector[string] {
	return func(r *expect.RequestInfo) string {
		return r.Request.Header.Get(name)
	}
}

func Query(name string) expect.ValueSelector[string] {
	return func(r *expect.RequestInfo) string {
		return r.Request.URL.Query().Get(name)
	}
}

func FormField(name string) expect.ValueSelector[string] {
	return func(r *expect.RequestInfo) string {
		return r.Request.Form.Get(name)
	}
}

func Body[T any](name string) expect.ValueSelector[T] {
	return func(r *expect.RequestInfo) T {
		return r.ParsedBody.(T)
	}
}
