package mocha

type MatcherParams struct {
	Config             Config
	Req                *MockRequest
	Repo               MockRepository
	ScenarioRepository ScenarioRepository
}

type Matcher[V any] func(v V, params MatcherParams) (bool, error)
