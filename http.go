package mocha

import (
	"github.com/vitorsalgado/mocha/v3/dzhttp"
	"github.com/vitorsalgado/mocha/v3/dzstd"
)

func NewAPI(config ...dzhttp.Configurer) *dzhttp.HTTPMockApp {
	return dzhttp.NewAPI(config...)
}

func NewAPIWithT(t dzstd.TestingT, config ...dzhttp.Configurer) *dzhttp.HTTPMockApp {
	return dzhttp.NewAPIWithT(t, config...)
}

func NewEchoServer(config ...dzhttp.Configurer) *dzhttp.HTTPMockApp {
	app := dzhttp.NewAPI(config...)
	app.MustMock(dzhttp.AnyMethod().Reply(dzhttp.Echo().Log()))

	return app
}
