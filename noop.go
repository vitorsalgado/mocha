package mocha

import (
	"github.com/vitorsalgado/mocha/core"
)

type noop struct {
}

func (n *noop) Cleanup(_ func()) {
}

func (n *noop) Errorf(_ string, _ ...any) {
}

func (n *noop) Helper() {
}

func Noop() core.T {
	return &noop{}
}
