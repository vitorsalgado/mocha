package mocha

import (
	"fmt"

	"github.com/vitorsalgado/mocha/core"
)

// StdoutNotifier implements core.T outputting logs to the stdout.
type StdoutNotifier struct {
}

func (n *StdoutNotifier) Logf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (n *StdoutNotifier) Errorf(format string, args ...any) {
	n.Logf(format, args...)
}

// FailNow do nothing.
func (n *StdoutNotifier) FailNow() {
}

// Helper do nothing.
func (n *StdoutNotifier) Helper() {
}

// NewStdoutNotifier returns a core.T implementation that logs to the stdout.
// FailNow() and Helper() will do nothing.
func NewStdoutNotifier() core.T {
	return &StdoutNotifier{}
}
