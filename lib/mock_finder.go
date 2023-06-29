package lib

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

// FindResult holds the results for an attempt to match a mock to a request.
type FindResult[MOCK Mock] struct {
	Pass            bool
	Matched         MOCK
	ClosestMatch    MOCK
	MismatchDetails []MismatchDetail
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
) *FindResult[MOCK] {
	var nearest MOCK
	var nearestPresent = false
	var aggWeight = 0
	var aggDetails = make([]MismatchDetail, 0)

	for _, m := range mocks {
		pass, weight, details := Match[TValueIn](requestValues, m.GetExpectations())
		if pass {
			return &FindResult[MOCK]{Pass: true, Matched: m}
		}

		if weight > 0 && weight > aggWeight {
			nearestPresent = true
			nearest = m
			aggWeight = weight
		}

		aggDetails = append(aggDetails, details...)
	}

	if nearestPresent {
		return &FindResult[MOCK]{Pass: false, ClosestMatch: nearest, MismatchDetails: aggDetails}
	}

	return &FindResult[MOCK]{Pass: false, MismatchDetails: aggDetails}
}

// Match checks if the current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func Match[VS any](ri VS, expectations []*Expectation[VS]) (bool, int, []MismatchDetail) {
	w := 0
	ok := true
	details := make([]MismatchDetail, 0, len(expectations))

	for _, exp := range expectations {
		var val any
		if exp.ValueSelector != nil {
			val = exp.ValueSelector(ri)
		}

		result, err := matchExpectation(exp, val)

		if err != nil {
			ok = false
			details = append(details, MismatchDetail{
				MatchersName: exp.Matcher.Name(),
				Target:       exp.Target,
				Key:          exp.Key,
				Err:          err,
			})

			continue
		}

		if result.Pass {
			w += int(exp.Weight)
		} else {
			ok = false
			details = append(details, MismatchDetail{
				MatchersName: exp.Matcher.Name(),
				Target:       exp.Target,
				Key:          exp.Key,
				Result:       result,
			})
		}
	}

	return ok, w, details
}

func matchExpectation[VS any](e *Expectation[VS], value any) (result *matcher.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: matcher=%s. %v", e.Matcher.Name(), r)
			return
		}
	}()

	result, err = e.Matcher.Match(value)
	if err != nil {
		err = fmt.Errorf("%s: error while matching. %w", e.Matcher.Name(), err)
	}

	return
}
