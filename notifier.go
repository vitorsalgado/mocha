package mocha

import (
	"fmt"

	"github.com/vitorsalgado/mocha/core"
)

type Notifier struct {
}

func (n *Notifier) Logf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (n *Notifier) Errorf(format string, args ...any) {
	n.Logf(format, args...)
}

func (n *Notifier) FailNow() {
}

func (n *Notifier) Helper() {
}

func StdoutNotifier() core.T {
	return &Notifier{}
}
