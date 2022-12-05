package matcher

import (
	"fmt"
	"sync/atomic"
)

type repeatMatcher struct {
	max  int64
	hits int64
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
	atomic.AddInt64(&m.hits, 1)
	return nil
}

func Repeat(times int64) Matcher {
	return &repeatMatcher{max: times}
}
