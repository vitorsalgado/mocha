package expect

const ScenarioStateStarted = "started"

type ScenarioState struct {
	Name  string
	State string
}

func NewScenarioState(name string) *ScenarioState {
	return &ScenarioState{Name: name, State: ScenarioStateStarted}
}

func (s *ScenarioState) HasStarted() bool {
	return s.State == ScenarioStateStarted
}

type ScenarioStorage interface {
	FetchByName(name string) (*ScenarioState, bool)
	CreateNewIfNeeded(name string) *ScenarioState
	Save(scenario *ScenarioState)
}

type internalScenarioStorage struct {
	data map[string]*ScenarioState
}

func NewScenarioStorage() ScenarioStorage {
	return &internalScenarioStorage{data: make(map[string]*ScenarioState)}
}

func (store *internalScenarioStorage) FetchByName(name string) (*ScenarioState, bool) {
	s, ok := store.data[name]
	return s, ok
}

func (store *internalScenarioStorage) CreateNewIfNeeded(name string) *ScenarioState {
	s, ok := store.FetchByName(name)

	if !ok {
		scenario := NewScenarioState(name)
		store.Save(scenario)
		return scenario
	}

	return s
}

func (store *internalScenarioStorage) Save(scenario *ScenarioState) {
	store.data[scenario.Name] = scenario
}

type ScenarioMatcher struct {
	Store         ScenarioStorage
	RequiredState string
	NewState      string
	Nm            string
}

func (m *ScenarioMatcher) Name() string {
	return "ScenarioState"
}

func (m *ScenarioMatcher) Match(_ any) (bool, error) {
	if m.RequiredState == ScenarioStateStarted {
		m.Store.CreateNewIfNeeded(m.Nm)
	}

	scn, ok := m.Store.FetchByName(m.Nm)
	if !ok {
		return true, nil
	}

	if scn.State == m.RequiredState {
		return true, nil
	}

	return false, nil
}

func (m *ScenarioMatcher) DescribeFailure(_ any) string {
	return ""
}

func (m *ScenarioMatcher) OnMockServed() error {
	scn, ok := m.Store.FetchByName(m.Nm)
	if !ok {
		return nil
	}

	if m.NewState != "" {
		scn.State = m.NewState
	}

	return nil
}

func Scenario(store ScenarioStorage) func(name, requiredState, newState string) Matcher {
	return func(name, requiredState, newState string) Matcher {
		return &ScenarioMatcher{
			Store:         store,
			RequiredState: requiredState,
			NewState:      newState,
			Nm:            name,
		}
	}
}
