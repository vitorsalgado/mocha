package mocha

type MatcherParams struct {
	Config        Config
	Request       *MockRequest
	MockStore     MockStore
	ScenarioStore ScenarioStore
}

type Matcher[V any] func(v V, params MatcherParams) (bool, error)
