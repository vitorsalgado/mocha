package matcher

import (
	"fmt"
	"sync/atomic"
)

type repeatMatcher struct {
	Max  int64
	Hits int64
}

func (m *repeatMatcher) Name() string {
	return "Repeat"
}

func (m *repeatMatcher) Match(_ any) (*Result, error) {
	return &Result{OK: m.Hits < m.Max, DescribeFailure: func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.Max)),
			_separator,
			printReceived(m.Hits),
		)
	}}, nil
}

func (m *repeatMatcher) OnMockServed() error {
	atomic.AddInt64(&m.Hits, 1)
	return nil
}

func Repeat(times int64) Matcher {
	return &repeatMatcher{Max: times}
}
