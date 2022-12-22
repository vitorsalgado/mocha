package logger

import (
	"fmt"
)

type Log interface {
	Logf(string, ...any)
}

var _ Log = (*Console)(nil)

// Console implements core.TestingT outputting logs to the stdout.
type Console struct {
}

func (n *Console) Logf(format string, args ...any) {
	fmt.Printf(format, args...)
}

// NewConsole returns a core.TestingT implementation that logs to the stdout.
// FailNow() and Helper() will do nothing.
func NewConsole() *Console {
	return &Console{}
}
