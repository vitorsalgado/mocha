package matcher

import (
	"fmt"
	"sync"
)

type timesMatcher struct {
	max  int
	hits int
	mu   sync.Mutex
}

func (m *timesMatcher) Name() string {
	return "Times"
}

func (m *timesMatcher) Match(_ any) (*Result, error) {
	return &Result{OK: m.hits < m.max, DescribeFailure: func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.max)),
			_separator,
			printReceived(m.hits),
		)
	}}, nil
}

func (m *timesMatcher) OnMockServed() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hits++

	return nil
}

func (m *timesMatcher) Spec() any {
	return []any{"times", m.max}
}

func Repeat(times int) Matcher {
	return &timesMatcher{max: times}
}
