package mocha

import (
	"github.com/vitorsalgado/mocha/v3/httpd"
	"github.com/vitorsalgado/mocha/v3/lib"
)

func NewAPI(config ...httpd.Configurer) *httpd.HTTPMockApp {
	return httpd.NewAPI(config...)
}

func NewAPIWithT(t lib.TestingT, config ...httpd.Configurer) *httpd.HTTPMockApp {
	return httpd.NewAPIWithT(t, config...)
}

func NewEchoServer(config ...httpd.Configurer) *httpd.HTTPMockApp {
	app := httpd.NewAPI(config...)
	app.MustMock(httpd.AnyMethod().Reply(httpd.Echo().Log()))

	return app
}
