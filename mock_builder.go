package mocha

import (
	"github.com/vitorsalgado/mocha/base"
	"net/http"
)

type MockBuilder struct {
	mock Mock
}

func NewBuilder() *MockBuilder {
	return &MockBuilder{mock: *NewMock()}
}

func (b *MockBuilder) Header(key string, matcher base.Matcher[string]) *MockBuilder {
	exp := buildExpectation(func(r *http.Request) string {
		return r.Header.Get(key)
	}, matcher)

	b.mock.Expectations = append(b.mock.Expectations, exp)

	return b
}

func (b *MockBuilder) Res() *MockBuilder {
	b.mock.ResFn = func(r *http.Request, mock *Mock) (Response, error) {
		return Response{Status: 201}, nil
	}

	return b
}

func (b *MockBuilder) Build() Mock {
	return b.mock
}

func buildExpectation[V any](picker RequestPicker[V], matcher base.Matcher[V]) Expectation[V] {
	return Expectation[V]{Pick: picker, Matcher: matcher}
}
