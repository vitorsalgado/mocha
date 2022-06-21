package mocha

import "github.com/vitorsalgado/mocha/matcher"

type findMockResult struct {
	Matches      bool
	Matched      *Mock
	ClosestMatch *Mock
}

func findMockForRequest(mockstore MockStore, params matcher.Params) (*findMockResult, error) {
	mocks := mockstore.FetchSorted()

	var m *Mock
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
