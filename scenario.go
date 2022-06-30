package mocha

import "github.com/vitorsalgado/mocha/matcher"

const (
	ScenarioStarted       = "STARTED"
	BuiltIntExtraScenario = "__mocha:scenarios"
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

	scenarioStore struct {
		data map[string]Scenario
	}
)

func NewScenarioStore() ScenarioStore {
	return &scenarioStore{data: make(map[string]Scenario)}
}

func (repo *scenarioStore) FetchByName(name string) *Scenario {
	s := repo.data[name]
	return &s
}

func (repo *scenarioStore) CreateNewIfNeeded(name string) *Scenario {
	s, ok := repo.data[name]

	if !ok {
		scenario := NewScenario(name)
		repo.Save(*scenario)
		return scenario
	}

	return &s
}

func (repo *scenarioStore) Save(scenario Scenario) {
	repo.data[scenario.Name] = scenario
}

func scenarioMatcher[V any](name, requiredState, newState string) matcher.Matcher[V] {
	return func(_ V, params matcher.Args) (bool, error) {
		s, _ := params.Params.Get(BuiltIntExtraScenario)
		scenarios := s.(ScenarioStore)

		if requiredState == ScenarioStarted {
			scenarios.CreateNewIfNeeded(name)
		}

		scenario := scenarios.FetchByName(name)

		if scenario == nil {
			return true, nil
		}

		if scenario.State == requiredState {
			if newState != "" {
				scenario.State = newState
				scenarios.Save(*scenario)
			}

			return true, nil
		}

		return false, nil
	}
}
