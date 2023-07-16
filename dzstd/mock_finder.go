package dzstd

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

// FindResult holds the results for an attempt to match a mock to a request.
type FindResult[MOCK Mock] struct {
	Pass    bool
	Matched MOCK

	// Closest matched mock.
	// Populated only in case of mismatch.
	ClosestMatch MOCK

	MismatchesCount int
}

type TMockMatcher[TValueIn any] interface {
	Mock
	RequestMatcher[TValueIn]
}

// FindMockForRequest tries to find a mock to the incoming HTTP request.
// It runs all matchers of all eligible mocks on request until it finds one that matches every one of them.
// It returns a FindResult with the find result, along with a possible closest match.
func FindMockForRequest[TValueIn any, MOCK TMockMatcher[TValueIn]](
	mocks []MOCK,
	requestValues TValueIn,
	desc *Description,
) *FindResult[MOCK] {
	var nearest MOCK
	var nearestPresent = false
	var aggWeight = 0
	var misses = 0

	for _, m := range mocks {
		pass, weight := Match[TValueIn](requestValues, desc, m.GetExpectations())
		if pass {
			return &FindResult[MOCK]{Pass: true, Matched: m}
		}

		if weight > 0 && weight > aggWeight {
			nearestPresent = true
			nearest = m
			aggWeight = weight
		}

		misses++
	}

	if nearestPresent {
		return &FindResult[MOCK]{Pass: false, ClosestMatch: nearest, MismatchesCount: misses}
	}

	return &FindResult[MOCK]{Pass: false, MismatchesCount: misses}
}

// Match checks if the current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func Match[VS any](ri VS, desc *Description, expectations []*Expectation[VS]) (bool, int) {
	passed, aggW := true, 0

	for i, exp := range expectations {
		var val any
		if exp.ValueSelector != nil {
			val = exp.ValueSelector(ri)
		}

		result, err := wrapMatch(exp, val, i)
		if err != nil {
			desc.AppendList(" ", exp.Target, err.Error())
			passed = false
			continue
		}

		if result.Pass {
			aggW += int(exp.Weight)
		} else {
			desc.AppendList(" ", exp.Target, result.Message)
			passed = false
		}
	}

	return passed, aggW
}

func wrapMatch[VS any](e *Expectation[VS], value any, idx int) (m matcher.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: matcher[%d]: %v", idx, r)
			return
		}
	}()

	return e.Matcher.Match(value)
}
