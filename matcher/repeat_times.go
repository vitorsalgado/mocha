package matcher

import (
	"fmt"
	"sync"

	"github.com/vitorsalgado/mocha/v3/types"
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
	return &Result{Pass: m.hits < m.max, Message: func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.max)),
			_separator,
			printReceived(m.hits),
		)
	}}, nil
}

func (m *timesMatcher) AfterMockSent() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hits++

	return nil
}

func (m *timesMatcher) Raw() types.RawValue {
	return types.RawValue{"times", m.max}
}

func Repeat(times int) Matcher {
	return &timesMatcher{max: times}
}
