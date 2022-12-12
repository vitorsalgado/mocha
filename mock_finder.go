package mocha

import (
	"github.com/vitorsalgado/mocha/v3/matcher"
)

// findResult holds the results for an attempt to match a mock to a request.
type findResult struct {
	Matches         bool
	Matched         *Mock
	ClosestMatch    *Mock
	MismatchDetails []mismatchDetail
}

// findMockForRequest tries to find a mock to the incoming HTTP request.
// It runs all matchers of all eligible mocks on request until it finds one that matches every one of then.
// It returns a findResult with the find result, along with a possible closest match.
func findMockForRequest(storage mockStore, ri *matcher.RequestInfo) *findResult {
	var mocks = storage.GetEligible()
	var matched *Mock
	var weights = 0
	var details = make([]mismatchDetail, 0)

	for _, m := range mocks {
		result := m.requestMatches(ri, m.expectations)

		if result.OK {
			return &findResult{Matches: true, Matched: m}
		}

		if result.Weight > 0 && result.Weight > weights {
			matched = m
			weights = result.Weight
		}

		details = append(details, result.MismatchDetails...)
	}

	if matched == nil {
		return &findResult{Matches: false, MismatchDetails: details}
	}

	return &findResult{Matches: false, ClosestMatch: matched, MismatchDetails: details}
}
