package matcher

import (
	"sync"
)

type repeatMatcher struct {
	max  int
	hits int
	mu   sync.Mutex
}

func (m *repeatMatcher) Name() string {
	return "Times"
}

func (m *repeatMatcher) Match(_ any) (*Result, error) {
	if m.hits < m.max {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{stringify(m.max)},
		Message: printReceived(m.hits),
	}, nil
}

func (m *repeatMatcher) AfterMockServed() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hits++

	return nil
}

func Repeat(times int) Matcher {
	return &repeatMatcher{max: times}
}
