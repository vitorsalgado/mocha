package matcher

import (
	"fmt"
)

const ScenarioStateStarted = "started"

type scenarioState struct {
	Name  string
	State string
}

func newScenarioState(name string) *scenarioState {
	return &scenarioState{Name: name, State: ScenarioStateStarted}
}

func (s *scenarioState) hasStarted() bool {
	return s.State == ScenarioStateStarted
}

type ScenarioStore interface {
	FetchByName(name string) (*scenarioState, bool)
	CreateNewIfNeeded(name string) *scenarioState
	Save(scenario *scenarioState)
}

type internalScenarioStorage struct {
	data map[string]*scenarioState
}

func NewScenarioStore() ScenarioStore {
	return &internalScenarioStorage{data: make(map[string]*scenarioState)}
}

func (store *internalScenarioStorage) FetchByName(name string) (*scenarioState, bool) {
	s, ok := store.data[name]
	return s, ok
}

func (store *internalScenarioStorage) CreateNewIfNeeded(name string) *scenarioState {
	s, ok := store.FetchByName(name)

	if !ok {
		scenario := newScenarioState(name)
		store.Save(scenario)
		return scenario
	}

	return s
}

func (store *internalScenarioStorage) Save(scenario *scenarioState) {
	store.data[scenario.Name] = scenario
}

type scenarioMatcher struct {
	Store         ScenarioStore
	RequiredState string
	NewState      string
	Nm            string
}

func (m *scenarioMatcher) Name() string {
	return "Scenario"
}

func (m *scenarioMatcher) Match(_ any) (*Result, error) {
	if m.RequiredState == ScenarioStateStarted {
		m.Store.CreateNewIfNeeded(m.Nm)
	}

	scn, ok := m.Store.FetchByName(m.Nm)
	if !ok {
		return &Result{OK: true}, nil
	}

	message := func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.RequiredState)),
			_separator,
			printReceived(scn.State),
		)
	}

	if scn.State == m.RequiredState {
		return &Result{OK: true}, nil
	}

	return &Result{OK: false, DescribeFailure: message}, nil
}

func (m *scenarioMatcher) OnMockServed() error {
	scn, ok := m.Store.FetchByName(m.Nm)
	if !ok {
		return nil
	}

	if m.NewState != "" {
		scn.State = m.NewState
	}

	return nil
}

func Scenario(store ScenarioStore) func(name, requiredState, newState string) Matcher {
	return func(name, requiredState, newState string) Matcher {
		return &scenarioMatcher{
			Store:         store,
			RequiredState: requiredState,
			NewState:      newState,
			Nm:            name,
		}
	}
}
