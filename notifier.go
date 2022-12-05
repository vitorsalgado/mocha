package mocha

import (
	"fmt"
)

// ConsoleNotifier implements core.TestingT outputting logs to the stdout.
type ConsoleNotifier struct {
}

func (n *ConsoleNotifier) Logf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (n *ConsoleNotifier) Errorf(format string, args ...any) {
	n.Logf(format, args...)
}

// FailNow do nothing.
func (n *ConsoleNotifier) FailNow() {
}

// Helper do nothing.
func (n *ConsoleNotifier) Helper() {
}

// NewConsoleNotifier returns a core.TestingT implementation that logs to the stdout.
// FailNow() and Helper() will do nothing.
func NewConsoleNotifier() TestingT {
	return &ConsoleNotifier{}
}
