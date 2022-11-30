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

type ScenarioStorage interface {
	FetchByName(name string) (*scenarioState, bool)
	CreateNewIfNeeded(name string) *scenarioState
	Save(scenario *scenarioState)
}

type internalScenarioStorage struct {
	data map[string]*scenarioState
}

func NewScenarioStorage() ScenarioStorage {
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

type ScenarioMatcher struct {
	Store         ScenarioStorage
	RequiredState string
	NewState      string
	Nm            string
}

func (m *ScenarioMatcher) Name() string {
	return "Scenario"
}

func (m *ScenarioMatcher) Match(_ any) (Result, error) {
	if m.RequiredState == ScenarioStateStarted {
		m.Store.CreateNewIfNeeded(m.Nm)
	}

	scn, ok := m.Store.FetchByName(m.Nm)
	if !ok {
		return Result{OK: true}, nil
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
		return Result{OK: true}, nil
	}

	return Result{OK: false, DescribeFailure: message}, nil
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
