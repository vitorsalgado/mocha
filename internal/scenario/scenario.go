package scenario

import "github.com/vitorsalgado/mocha/matcher"

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

func Scenario[V any](name, requiredState, newState string) matcher.Matcher[V] {
	return func(_ V, params matcher.Args) (bool, error) {
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
}
