package notifier

import (
	"fmt"
)

// Console implements core.TestingT outputting logs to the stdout.
type Console struct {
}

func (n *Console) Logf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (n *Console) Errorf(format string, args ...any) {
	n.Logf(format, args...)
}

// Helper do nothing.
func (n *Console) Helper() {
}

// NewConsole returns a core.TestingT implementation that logs to the stdout.
// FailNow() and Helper() will do nothing.
func NewConsole() *Console {
	return &Console{}
}
