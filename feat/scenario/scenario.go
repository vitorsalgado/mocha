// Package scenario implements a stateful matcher that works like a stateful machine were states can be assigned during
// mock configuration.
// Mocks can be configured to be returned on certain state values and also,
// mocks can define new state values on they are served.
package scenario

import "github.com/vitorsalgado/mocha/expect"

const (
	StateStarted      = "STARTED"
	BuiltInParamStore = "__mocha:scenarios"
)

type scenario struct {
	Name  string
	State string
}

func newScenario(name string) scenario {
	return scenario{Name: name, State: StateStarted}
}

// HasStarted returns true when Scenario state is equal to "STARTED"
func (s scenario) HasStarted() bool {
	return s.State == StateStarted
}

type (
	Store interface {
		FetchByName(name string) (scenario, bool)
		CreateNewIfNeeded(name string) scenario
		Save(scenario scenario)
	}

	scenarioStore struct {
		data map[string]scenario
	}
)

func NewStore() Store {
	return &scenarioStore{data: make(map[string]scenario)}
}

func (store *scenarioStore) FetchByName(name string) (scenario, bool) {
	s, ok := store.data[name]
	return s, ok
}

func (store *scenarioStore) CreateNewIfNeeded(name string) scenario {
	s, ok := store.FetchByName(name)

	if !ok {
		scenario := newScenario(name)
		store.Save(scenario)
		return scenario
	}

	return s
}

func (store *scenarioStore) Save(scenario scenario) {
	store.data[scenario.Name] = scenario
}

func Scenario(name, requiredState, newState string) expect.Matcher {
	m := expect.Matcher{}
	m.Name = "Scenario"
	m.Matches = func(_ any, params expect.Args) (bool, error) {
		s, _ := params.Params.Get(BuiltInParamStore)
		scenarios := s.(Store)

		if requiredState == StateStarted {
			scenarios.CreateNewIfNeeded(name)
		}

		scn, ok := scenarios.FetchByName(name)
		if !ok {
			return true, nil
		}

		if scn.State == requiredState {
			if newState != "" {
				scn.State = newState
				scenarios.Save(scn)
			}

			return true, nil
		}

		return false, nil
	}

	return m
}
