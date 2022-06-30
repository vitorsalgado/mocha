package mocha

import (
	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/mock"
)

type findMockResult struct {
	Matches      bool
	Matched      *mock.Mock
	ClosestMatch *mock.Mock
}

func findMockForRequest(storage mock.Storage, params matcher.Args) (*findMockResult, error) {
	var mocks = storage.FetchEligible()
	var matched *mock.Mock
	var w = 0

	for _, m := range mocks {
		matches, err := m.Matches(params)
		if err != nil {
			return nil, err
		}

		if matches.IsMatch {
			return &findMockResult{Matches: true, Matched: m}, nil
		}

		if matches.Weight > 0 && matches.Weight > w {
			matched = m
			w = matches.Weight
		}
	}

	if matched == nil {
		return &findMockResult{Matches: false}, nil
	}

	return &findMockResult{Matches: false, ClosestMatch: matched}, nil
}
