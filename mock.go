package mocha

import (
	"github.com/vitorsalgado/mocha/base"
	"net/http"
)

type (
	Mock struct {
		ID           int32
		Name         string
		Priority     int
		Expectations []any
		ResFn        ResponseDelegate
		Hits         int
	}

	RequestPicker[V any] func(r *http.Request) V

	Expectation[V any] struct {
		Matcher base.Matcher[V]
		Pick    RequestPicker[V]
	}
)

func NewMock() *Mock {
	return &Mock{}
}
