package expect

import (
	"sync/atomic"
)

type RepeatMatcher struct {
	Times int64
	Hits  int64
}

func (m *RepeatMatcher) Name() string {
	return "Repeat"
}

func (m *RepeatMatcher) Match(_ any) (bool, error) {
	return m.Hits < m.Times, nil
}

func (m *RepeatMatcher) DescribeFailure(v any) string {
	return ""
}

func (m *RepeatMatcher) OnMockServed() error {
	atomic.AddInt64(&m.Hits, 1)
	return nil
}

func Repeat(times int64) Matcher {
	return &RepeatMatcher{Times: times}
}
