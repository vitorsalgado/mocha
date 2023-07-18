package mocha

import (
	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

func NewAPI(config ...dzhttp.Configurer) *dzhttp.HTTPMockApp {
	return dzhttp.NewAPI(config...)
}

func NewEchoServer(config ...dzhttp.Configurer) *dzhttp.HTTPMockApp {
	app := dzhttp.NewAPI(config...)
	app.MustMock(dzhttp.AnyMethod().Reply(dzhttp.Echo().Log()))

	return app
}
