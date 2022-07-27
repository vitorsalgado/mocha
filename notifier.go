package mocha

import (
	"fmt"
)

// T is based on testing.T and allow mocha components to log information and errors.
type T interface {
	Helper()
	Logf(string, ...any)
	Errorf(string, ...any)
	FailNow()
}

// ConsoleNotifier implements core.T outputting logs to the stdout.
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

// NewConsoleNotifier returns a core.T implementation that logs to the stdout.
// FailNow() and Helper() will do nothing.
func NewConsoleNotifier() T {
	return &ConsoleNotifier{}
}
