package matcher

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/types"
)

const ScenarioStateStarted = "started"

type scenarioState struct {
	name  string
	state string
}

func newScenarioState(name string) *scenarioState {
	return &scenarioState{name: name, state: ScenarioStateStarted}
}

func (s *scenarioState) hasStarted() bool {
	return s.state == ScenarioStateStarted
}

type scenarioStore struct {
	data map[string]*scenarioState
}

func newScenarioStore() *scenarioStore {
	return &scenarioStore{data: make(map[string]*scenarioState)}
}

func (store *scenarioStore) fetchByName(name string) (*scenarioState, bool) {
	s, ok := store.data[name]
	return s, ok
}

func (store *scenarioStore) createNewIfNeeded(name string) *scenarioState {
	s, ok := store.fetchByName(name)

	if !ok {
		scenario := newScenarioState(name)
		store.save(scenario)
		return scenario
	}

	return s
}

func (store *scenarioStore) save(scenario *scenarioState) {
	store.data[scenario.name] = scenario
}

type scenarioMatcher struct {
	store         *scenarioStore
	requiredState string
	newState      string
	nm            string
}

func (m *scenarioMatcher) Name() string {
	return "Scenario"
}

func (m *scenarioMatcher) Match(_ any) (*Result, error) {
	if m.requiredState == ScenarioStateStarted {
		m.store.createNewIfNeeded(m.nm)
	}

	scn, ok := m.store.fetchByName(m.nm)
	if !ok {
		return &Result{Pass: true}, nil
	}

	message := func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.requiredState)),
			_separator,
			printReceived(scn.state),
		)
	}

	if scn.state == m.requiredState {
		return &Result{Pass: true}, nil
	}

	return &Result{Pass: false, Message: message}, nil
}

func (m *scenarioMatcher) AfterMockSent() error {
	scn, ok := m.store.fetchByName(m.nm)
	if !ok {
		return nil
	}

	if m.newState != "" {
		scn.state = m.newState
	}

	return nil
}

func (m *scenarioMatcher) Raw() types.RawValue {
	return nil
}

func Scenario(name, requiredState, newState string) Matcher {
	return &scenarioMatcher{
		store:         newScenarioStore(),
		requiredState: requiredState,
		newState:      newState,
		nm:            name,
	}
}
