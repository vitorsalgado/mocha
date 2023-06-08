package lib

import (
	"reflect"
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
	var matched MOCK
	var weights = 0
	var details = make([]MismatchDetail, 0)

	for _, m := range mocks {
		result := Match[TValueIn](requestValues, m.GetExpectations())

		if result.Pass {
			return &FindResult[MOCK]{Pass: true, Matched: m}
		}

		if result.Weight > 0 && result.Weight > weights {
			matched = m
			weights = result.Weight
		}

		details = append(details, result.Details...)
	}

	if reflect.ValueOf(matched).IsNil() {
		return &FindResult[MOCK]{Pass: false, MismatchDetails: details}
	}

	return &FindResult[MOCK]{Pass: false, ClosestMatch: matched, MismatchDetails: details}
}

// Match checks if the current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func Match[VS any](ri VS, expectations []*Expectation[VS]) *MatchResult {
	w := 0
	ok := true
	details := make([]MismatchDetail, 0)

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

	return &MatchResult{Pass: ok, Weight: w, Details: details}
}
