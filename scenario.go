package mocha

import "github.com/vitorsalgado/mocha/matcher"

const (
	ScenarioStarted      = "STARTED"
	BuiltInParamScenario = "__mocha:scenarios"
)

type Scenario struct {
	Name  string
	State string
}

func NewScenario(name string) Scenario {
	return Scenario{Name: name, State: ScenarioStarted}
}

func (s Scenario) HasStarted() bool {
	return s.State == ScenarioStarted
}

type (
	ScenarioStore interface {
		FetchByName(name string) (Scenario, bool)
		CreateNewIfNeeded(name string) Scenario
		Save(scenario Scenario)
	}

	scenarioStore struct {
		data map[string]Scenario
	}
)

func NewScenarioStore() ScenarioStore {
	return &scenarioStore{data: make(map[string]Scenario)}
}

func (store *scenarioStore) FetchByName(name string) (Scenario, bool) {
	s, ok := store.data[name]
	return s, ok
}

func (store *scenarioStore) CreateNewIfNeeded(name string) Scenario {
	s, ok := store.FetchByName(name)

	if !ok {
		scenario := NewScenario(name)
		store.Save(scenario)
		return scenario
	}

	return s
}

func (store *scenarioStore) Save(scenario Scenario) {
	store.data[scenario.Name] = scenario
}

func scenarioMatcher[V any](name, requiredState, newState string) matcher.Matcher[V] {
	return func(_ V, params matcher.Args) (bool, error) {
		s, _ := params.Params.Get(BuiltInParamScenario)
		scenarios := s.(ScenarioStore)

		if requiredState == ScenarioStarted {
			scenarios.CreateNewIfNeeded(name)
		}

		scenario, ok := scenarios.FetchByName(name)
		if !ok {
			return true, nil
		}

		if scenario.State == requiredState {
			if newState != "" {
				scenario.State = newState
				scenarios.Save(scenario)
			}

			return true, nil
		}

		return false, nil
	}
}
