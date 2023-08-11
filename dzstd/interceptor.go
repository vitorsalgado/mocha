package dzstd

import "io"

type Chain interface {
	Next([]byte) (n int, err error)
}

type InterceptorChain struct {
	interceptors []Interceptor
	index        int
}

func NewChain(interceptors []Interceptor) Chain {
	return &InterceptorChain{interceptors: interceptors}
}

func (in *InterceptorChain) Next(p []byte) (n int, err error) {
	next := InterceptorChain{in.interceptors, in.index + 1}
	intereptor := in.interceptors[in.index]

	return intereptor.Intercept(p, &next)
}

type Interceptor interface {
	Intercept(p []byte, chain Chain) (n int, err error)
}

type RootIntereptor struct{ W io.Writer }

func (r *RootIntereptor) Intercept(p []byte, _ Chain) (n int, err error) {
	return r.W.Write(p)
}
