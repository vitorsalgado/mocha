package mocha

import (
	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/httpd"
)

func NewAPI(config ...mhttp.Configurer) *mhttp.HTTPMockApp {
	return mhttp.NewAPI(config...)
}

func NewAPIWithT(t foundation.TestingT, config ...mhttp.Configurer) *mhttp.HTTPMockApp {
	return mhttp.NewAPIWithT(t, config...)
}
