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
	ScenarioRepository interface {
		FetchByName(name string) *Scenario
		CreateNewIfNeeded(name string) *Scenario
		Save(scenario Scenario)
	}

	scenarioInMemoryRepository struct {
		data map[string]Scenario
	}
)

func NewScenarioRepository() ScenarioRepository {
	return &scenarioInMemoryRepository{data: make(map[string]Scenario)}
}

func (repo *scenarioInMemoryRepository) FetchByName(name string) *Scenario {
	s := repo.data[name]
	return &s
}

func (repo *scenarioInMemoryRepository) CreateNewIfNeeded(name string) *Scenario {
	s, ok := repo.data[name]

	if !ok {
		scenario := NewScenario(name)
		repo.Save(*scenario)
		return scenario
	}

	return &s
}

func (repo *scenarioInMemoryRepository) Save(scenario Scenario) {
	repo.data[scenario.Name] = scenario
}

func scenarioMatcher[V any](name, requiredState, newState string) Matcher[V] {
	return func(_ V, params MatcherParams) (bool, error) {
		if requiredState == ScenarioStarted {
			params.ScenarioRepository.CreateNewIfNeeded(name)
		}

		scenario := params.ScenarioRepository.FetchByName(name)

		if scenario == nil {
			return true, nil
		}

		if scenario.State == requiredState {
			if newState != "" {
				scenario.State = newState
				params.ScenarioRepository.Save(*scenario)
			}

			return true, nil
		}

		return false, nil
	}
}
