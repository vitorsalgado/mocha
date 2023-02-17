package matcher

import (
	"sync"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
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
		Ext:     []string{mfmt.Stringify(m.max)},
		Message: mfmt.PrintReceived(m.hits),
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
