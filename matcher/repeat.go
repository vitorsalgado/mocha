package matcher

import (
	"fmt"
	"sync"
)

type repeatMatcher struct {
	max  int
	hits int
	mu   sync.Mutex
}

func (m *repeatMatcher) Name() string {
	return "Repeat"
}

func (m *repeatMatcher) Match(_ any) (*Result, error) {
	return &Result{OK: m.hits < m.max, DescribeFailure: func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.max)),
			_separator,
			printReceived(m.hits),
		)
	}}, nil
}

func (m *repeatMatcher) OnMockServed() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hits++

	return nil
}

func Repeat(times int) Matcher {
	return &repeatMatcher{max: times}
}
