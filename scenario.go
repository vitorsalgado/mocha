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

	ScenarioInMemoryRepository struct {
		data map[string]Scenario
	}
)

func NewScenarioRepository() ScenarioRepository {
	return &ScenarioInMemoryRepository{data: make(map[string]Scenario)}
}

func (repo *ScenarioInMemoryRepository) FetchByName(name string) *Scenario {
	s := repo.data[name]
	return &s
}

func (repo *ScenarioInMemoryRepository) CreateNewIfNeeded(name string) *Scenario {
	s, ok := repo.data[name]

	if !ok {
		scenario := NewScenario(name)
		repo.Save(*scenario)
		return scenario
	}

	return &s
}

func (repo *ScenarioInMemoryRepository) Save(scenario Scenario) {
	repo.data[scenario.Name] = scenario
}

func ScenarioM[V any](name, requiredState, newState string) Matcher[V] {
	return func(_ V, ctx MatcherContext) (bool, error) {
		if requiredState == ScenarioStarted {
			ctx.ScenarioRepository.CreateNewIfNeeded(name)
		}

		scenario := ctx.ScenarioRepository.FetchByName(name)
		if scenario != nil {
			if scenario.State == requiredState {
				if newState != "" {
					scenario.State = newState
					ctx.ScenarioRepository.Save(*scenario)
				}

				return true, nil
			} else {
				return false, nil
			}
		}

		return true, nil
	}
}
