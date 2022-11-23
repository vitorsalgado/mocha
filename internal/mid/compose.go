package mid

import "net/http"

type Middlewares struct {
	handlers []func(handler http.Handler) http.Handler
}

func Compose(handlers ...func(handler http.Handler) http.Handler) Middlewares {
	return Middlewares{handlers: append(([]func(handler http.Handler) http.Handler)(nil), handlers...)}
}

func (m Middlewares) Root(root http.Handler) http.Handler {
	for i := range m.handlers {
		root = m.handlers[len(m.handlers)-1-i](root)
	}

	return root
}
