package mocha

import (
	"github.com/vitorsalgado/mocha/v2/expect"
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
func findMockForRequest(storage storage, params expect.Args) (*findResult, error) {
	var mocks = storage.FetchEligible()
	var matched *Mock
	var weights = 0
	var details = make([]mismatchDetail, 0)

	for _, m := range mocks {
		result, err := m.matches(params, m.Expectations)
		if err != nil {
			return nil, err
		}

		if result.IsMatch {
			return &findResult{Matches: true, Matched: m}, nil
		}

		if result.Weight > 0 && result.Weight > weights {
			matched = m
			weights = result.Weight
		}

		details = append(details, result.MismatchDetails...)
	}

	if matched == nil {
		return &findResult{Matches: false, MismatchDetails: details}, nil
	}

	return &findResult{Matches: false, ClosestMatch: matched, MismatchDetails: details}, nil
}
