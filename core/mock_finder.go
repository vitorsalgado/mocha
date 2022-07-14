package core

import (
	"github.com/vitorsalgado/mocha/expect"
)

// FindResult holds the results for an attempt to match a mock to a request.
type FindResult struct {
	Matches         bool
	Matched         *Mock
	ClosestMatch    *Mock
	MismatchDetails []MismatchDetail
}

// FindMockForRequest tries to find a mock to the incoming HTTP request.
// It runs all matchers of all eligible mocks on request until it finds one that matches every one of then.
// It returns a FindResult with the find result, along with a possible closest match.
func FindMockForRequest(storage Storage, params expect.Args) (*FindResult, error) {
	var mocks = storage.FetchEligible()
	var matched *Mock
	var weights = 0
	var details = make([]MismatchDetail, 0)

	for _, m := range mocks {
		result, err := m.Matches(params, m.Expectations)
		if err != nil {
			return nil, err
		}

		if result.IsMatch {
			return &FindResult{Matches: true, Matched: m}, nil
		}

		if result.Weight > 0 && result.Weight > weights {
			matched = m
			weights = result.Weight
		}

		details = append(details, result.MismatchDetails...)
	}

	if matched == nil {
		return &FindResult{Matches: false, MismatchDetails: details}, nil
	}

	return &FindResult{Matches: false, ClosestMatch: matched, MismatchDetails: details}, nil
}
