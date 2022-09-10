package mocha

const (
	_scenarioStateStarted = "STARTED"
)

type scenario struct {
	Name  string
	State string
}

func newScenario(name string) scenario {
	return scenario{Name: name, State: _scenarioStateStarted}
}

// HasStarted returns true when Scenario state is equal to "STARTED"
func (s scenario) HasStarted() bool {
	return s.State == _scenarioStateStarted
}

type (
	scenarioStore interface {
		FetchByName(name string) (scenario, bool)
		CreateNewIfNeeded(name string) scenario
		Save(s scenario)
	}

	internalScenarioStore struct {
		data map[string]scenario
	}
)

func newScenarioStore() scenarioStore {
	return &internalScenarioStore{data: make(map[string]scenario)}
}

func (store *internalScenarioStore) FetchByName(name string) (scenario, bool) {
	s, ok := store.data[name]
	return s, ok
}

func (store *internalScenarioStore) CreateNewIfNeeded(name string) scenario {
	s, ok := store.FetchByName(name)

	if !ok {
		sc := newScenario(name)
		store.Save(sc)
		return sc
	}

	return s
}

func (store *internalScenarioStore) Save(s scenario) {
	store.data[s.Name] = s
}
