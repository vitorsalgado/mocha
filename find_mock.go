package mocha

type FindMockResult struct {
	Matches      bool
	Matched      *Mock
	ClosestMatch *Mock
}

func FindMockForRequest(params MatcherParams) (*FindMockResult, error) {
	mocks := params.Repo.FetchSorted()

	var m *Mock
	var w = 0

	for _, mock := range mocks {
		matches, err := mock.Matches(params)
		if err != nil {
			return nil, err
		}

		if matches.IsMatch {
			return &FindMockResult{Matches: true, Matched: &mock}, nil
		}

		if matches.Weight > 0 && matches.Weight > w {
			m = &mock
			w = matches.Weight
		}
	}

	if m == nil {
		return &FindMockResult{Matches: false}, nil
	}

	return &FindMockResult{Matches: false, ClosestMatch: m}, nil
}
