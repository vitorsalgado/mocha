package middleware

import "net/http"

type Handler func(handler http.Handler) http.Handler

type Middlewares struct {
	handlers []Handler
}

func Compose(handlers ...Handler) Middlewares {
	return Middlewares{handlers: append(([]Handler)(nil), handlers...)}
}

func (m Middlewares) Root(root http.Handler) http.Handler {
	for i := range m.handlers {
		root = m.handlers[len(m.handlers)-1-i](root)
	}

	return root
}
