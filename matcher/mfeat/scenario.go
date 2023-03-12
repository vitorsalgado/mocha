package mfeat

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher"
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

type ScenarioStore struct {
	data map[string]*scenarioState
}

func NewScenarioStore() *ScenarioStore {
	return &ScenarioStore{data: make(map[string]*scenarioState)}
}

func (store *ScenarioStore) fetchByName(name string) (*scenarioState, bool) {
	s, ok := store.data[strings.ToLower(strings.TrimSpace(name))]
	return s, ok
}

func (store *ScenarioStore) createNewIfNeeded(name string) *scenarioState {
	s, ok := store.fetchByName(strings.ToLower(strings.TrimSpace(name)))

	if !ok {
		scenario := newScenarioState(name)
		store.save(scenario)
		return scenario
	}

	return s
}

func (store *ScenarioStore) save(scenario *scenarioState) {
	store.data[scenario.name] = scenario
}

type scenarioMatcher struct {
	store         *ScenarioStore
	requiredState string
	newState      string
	name          string
}

func (m *scenarioMatcher) Name() string {
	return "Scenario"
}

func (m *scenarioMatcher) Match(_ any) (*matcher.Result, error) {
	if m.requiredState == ScenarioStateStarted {
		m.store.createNewIfNeeded(m.name)
	}

	scn, ok := m.store.fetchByName(m.name)
	if !ok {
		return &matcher.Result{Pass: true}, nil
	}

	if scn.state == m.requiredState {
		return &matcher.Result{Pass: true}, nil
	}

	return &matcher.Result{
		Ext:     []string{m.name, scn.state},
		Message: fmt.Sprintf("required scenario state: %s", m.requiredState),
	}, nil
}

func (m *scenarioMatcher) AfterMockServed() error {
	scn, ok := m.store.fetchByName(m.name)
	if !ok {
		return nil
	}

	if m.newState != "" {
		scn.state = m.newState
	}

	return nil
}

func Scenario(store *ScenarioStore, name, requiredState, newState string) matcher.Matcher {
	return &scenarioMatcher{
		store:         store,
		requiredState: requiredState,
		newState:      newState,
		name:          name,
	}
}
