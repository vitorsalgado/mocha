package mock

import (
	"github.com/vitorsalgado/mocha/to"
)

// FindResult holds the results for an attempt to match a mock to a request.
type FindResult struct {
	Matches      bool
	Matched      *Mock
	ClosestMatch *Mock
}

// FindForRequest tries to find a mock to the incoming HTTP request.
// It runs all matchers of all eligible mocks on request until it finds one that matches every one of then.
// It returns a FindResult with the find result, along with a possible closest match.
func FindForRequest(storage Storage, params to.Args, t T) (*FindResult, error) {
	var mocks = storage.FetchEligible()
	var matched *Mock
	var w = 0

	for _, m := range mocks {
		result, err := m.Matches(params, m.Expectations, t)
		if err != nil {
			return nil, err
		}

		if result.IsMatch {
			return &FindResult{Matches: true, Matched: m}, nil
		}

		if result.Weight > 0 && result.Weight > w {
			matched = m
			w = result.Weight
		}
	}

	if matched == nil {
		return &FindResult{Matches: false}, nil
	}

	return &FindResult{Matches: false, ClosestMatch: matched}, nil
}
