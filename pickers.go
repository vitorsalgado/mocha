package mocha

import (
	"github.com/vitorsalgado/mocha/expect"
)

func Header(name string) expect.ValueSelector[string] {
	return func(r *expect.RequestInfo) string {
		return r.Request.Header.Get(name)
	}
}
