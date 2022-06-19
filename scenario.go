package mocha

const (
	ScenarioStarted = "STARTED"
)

type Scenario struct {
	Name  string
	State string
}

func NewScenario(name string) *Scenario {
	return &Scenario{Name: name, State: ScenarioStarted}
}

func (s Scenario) HasStarted() bool {
	return s.State == ScenarioStarted
}

type (
	ScenarioStore interface {
		FetchByName(name string) *Scenario
		CreateNewIfNeeded(name string) *Scenario
		Save(scenario Scenario)
	}

	scenarioInMemoStore struct {
		data map[string]Scenario
	}
)

func NewScenarioStore() ScenarioStore {
	return &scenarioInMemoStore{data: make(map[string]Scenario)}
}

func (repo *scenarioInMemoStore) FetchByName(name string) *Scenario {
	s := repo.data[name]
	return &s
}

func (repo *scenarioInMemoStore) CreateNewIfNeeded(name string) *Scenario {
	s, ok := repo.data[name]

	if !ok {
		scenario := NewScenario(name)
		repo.Save(*scenario)
		return scenario
	}

	return &s
}

func (repo *scenarioInMemoStore) Save(scenario Scenario) {
	repo.data[scenario.Name] = scenario
}

func scenarioMatcher[V any](name, requiredState, newState string) Matcher[V] {
	return func(_ V, params MatcherParams) (bool, error) {
		if requiredState == ScenarioStarted {
			params.ScenarioStore.CreateNewIfNeeded(name)
		}

		scenario := params.ScenarioStore.FetchByName(name)

		if scenario == nil {
			return true, nil
		}

		if scenario.State == requiredState {
			if newState != "" {
				scenario.State = newState
				params.ScenarioStore.Save(*scenario)
			}

			return true, nil
		}

		return false, nil
	}
}
