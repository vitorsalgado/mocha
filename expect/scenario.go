package expect

const (
	ScenarioStateStarted      = "STARTED"
	ScenarioBuiltInParamStore = "@@mocha:scenarios"
)

type scenario struct {
	Name  string
	State string
}

func newScenario(name string) scenario {
	return scenario{Name: name, State: ScenarioStateStarted}
}

// HasStarted returns true when Scenario state is equal to "STARTED"
func (s scenario) HasStarted() bool {
	return s.State == ScenarioStateStarted
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

func NewScenarioStore() Store {
	return &scenarioStore{data: make(map[string]scenario)}
}

func (store *scenarioStore) FetchByName(name string) (scenario, bool) {
	s, ok := store.data[name]
	return s, ok
}

func (store *scenarioStore) CreateNewIfNeeded(name string) scenario {
	s, ok := store.FetchByName(name)

	if !ok {
		sc := newScenario(name)
		store.Save(sc)
		return sc
	}

	return s
}

func (store *scenarioStore) Save(scenario scenario) {
	store.data[scenario.Name] = scenario
}

func Scenario(name, requiredState, newState string) Matcher {
	m := Matcher{}
	m.Name = "Scenario"
	m.Matches = func(_ any, params Args) (bool, error) {
		s, _ := params.Params.Get(ScenarioBuiltInParamStore)
		scenarios := s.(Store)

		if requiredState == ScenarioStateStarted {
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
