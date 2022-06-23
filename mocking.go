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

func findMockForRequest(mockstore mock.Storage, params matcher.Params) (*findMockResult, error) {
	mocks := mockstore.FetchSorted()

	var m *mock.Mock
	var w = 0

	for _, mock := range mocks {
		matches, err := mock.Matches(params)
		if err != nil {
			return nil, err
		}

		if matches.IsMatch {
			return &findMockResult{Matches: true, Matched: &mock}, nil
		}

		if matches.Weight > 0 && matches.Weight > w {
			m = &mock
			w = matches.Weight
		}
	}

	if m == nil {
		return &findMockResult{Matches: false}, nil
	}

	return &findMockResult{Matches: false, ClosestMatch: m}, nil
}
