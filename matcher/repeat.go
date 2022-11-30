package matcher

import (
	"fmt"
	"sync/atomic"
)

type RepeatMatcher struct {
	Max  int64
	Hits int64
}

func (m *RepeatMatcher) Name() string {
	return "Repeat"
}

func (m *RepeatMatcher) Match(_ any) (Result, error) {
	return Result{OK: m.Hits < m.Max, DescribeFailure: func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.Max)),
			_separator,
			printReceived(m.Hits),
		)
	}}, nil
}

func (m *RepeatMatcher) OnMockServed() error {
	atomic.AddInt64(&m.Hits, 1)
	return nil
}

func Repeat(times int64) Matcher {
	return &RepeatMatcher{Max: times}
}
