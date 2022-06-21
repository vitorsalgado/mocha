package mocha

import (
	"fmt"

	"github.com/vitorsalgado/mocha/matcher"
)

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
	ScenarioStore struct {
		data map[string]Scenario
	}
)

func NewScenarioStore() *ScenarioStore {
	return &ScenarioStore{data: make(map[string]Scenario)}
}

func (repo *ScenarioStore) FetchByName(name string) *Scenario {
	s := repo.data[name]
	return &s
}

func (repo *ScenarioStore) CreateNewIfNeeded(name string) *Scenario {
	s, ok := repo.data[name]

	if !ok {
		scenario := NewScenario(name)
		repo.Save(*scenario)
		return scenario
	}

	return &s
}

func (repo *ScenarioStore) Save(scenario Scenario) {
	repo.data[scenario.Name] = scenario
}

func scenarioMatcher[V any](name, requiredState, newState string) matcher.Matcher[V] {
	return func(_ V, params matcher.Params) (bool, error) {
		s, _ := params.Extras.Get("scenarios")
		scenarios, e := s.(*ScenarioStore)
		if !e {
			return false, fmt.Errorf("")
		}

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
