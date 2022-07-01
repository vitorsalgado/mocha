package matcher

const (
	ScenarioStarted      = "STARTED"
	BuiltInParamScenario = "__mocha:scenarios"
)

type scenario struct {
	Name  string
	State string
}

func newScenario(name string) scenario {
	return scenario{Name: name, State: ScenarioStarted}
}

func (s scenario) HasStarted() bool {
	return s.State == ScenarioStarted
}

type (
	ScenarioStore interface {
		FetchByName(name string) (scenario, bool)
		CreateNewIfNeeded(name string) scenario
		Save(scenario scenario)
	}

	scenarioStore struct {
		data map[string]scenario
	}
)

func NewScenarioStore() ScenarioStore {
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

func Scenario[V any](name, requiredState, newState string) Matcher[V] {
	return func(_ V, params Args) (bool, error) {
		s, _ := params.Params.Get(BuiltInParamScenario)
		scenarios := s.(ScenarioStore)

		if requiredState == ScenarioStarted {
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
