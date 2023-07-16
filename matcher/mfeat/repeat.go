package mfeat

import (
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

type repeatMatcher struct {
	max  int64
	hits int64
}

func (m *repeatMatcher) Match(_ any) (matcher.Result, error) {
	if atomic.LoadInt64(&m.hits) < m.max {
		return matcher.Result{Pass: true}, nil
	}

	return matcher.Result{
		Message: strings.Join([]string{"Repeat(", strconv.FormatInt(m.max, 10), ") Reached the max matched requests for this mock"}, ""),
	}, nil
}

func (m *repeatMatcher) AfterMockServed() error {
	atomic.AddInt64(&m.hits, 1)

	return nil
}

func Repeat(times int64) matcher.Matcher {
	return &repeatMatcher{max: times}
}
