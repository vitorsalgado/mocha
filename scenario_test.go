package mocha

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/reply"
)

func TestScenario(t *testing.T) {
	t.Run("should init scenario as started", func(t *testing.T) {
		assert.True(t, NewScenario("test").HasStarted())
	})

	t.Run("should only create scenario if needed", func(t *testing.T) {
		store := NewScenarioStore()
		store.CreateNewIfNeeded("scenario-1")

		s, ok := store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.True(t, s.HasStarted())

		s.State = "another-state"
		store.Save(s)

		store.CreateNewIfNeeded("scenario-1")

		s, ok = store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.False(t, s.HasStarted())
		assert.Equal(t, s.State, "another-state")
	})
}

func TestScenarioConditions(t *testing.T) {
	store := NewScenarioStore()
	params := params.New()
	params.Set(BuiltInParamScenario, store)
	args := matcher.Args{Params: params}

	t.Run("should return true when scenario is not started and also not found", func(t *testing.T) {
		m := scenarioMatcher[any]("test", "required", "new")
		res, err := m(nil, args)

		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return false when scenario exists but it is not in the required state", func(t *testing.T) {
		store.CreateNewIfNeeded("hi")

		m := scenarioMatcher[any]("hi", "required", "new")
		res, err := m(nil, args)

		assert.Nil(t, err)
		assert.False(t, res)
	})
}

func TestScenarioMatcher(t *testing.T) {
	m := ForTest(t)
	m.Start()
	scenario := "test"

	s1 := m.Mock(Get(matcher.URLPath("/1")).
		StartScenario(scenario).
		ScenarioStateWillBe("step2").
		Name("step-1").
		Reply(reply.OK().BodyString("step1")))

	s2 := m.Mock(Get(matcher.URLPath("/2")).
		ScenarioIs(scenario).
		ScenarioStateIs("step2").
		ScenarioStateWillBe("step3").
		Name("step-2").
		Reply(reply.OK().BodyString("step2")))

	s3 := m.Mock(Get(matcher.URLPath("/3")).
		ScenarioIs(scenario).
		ScenarioStateIs("step3").
		ScenarioStateWillBe("step4").
		Name("step-3").
		Reply(reply.OK().BodyString("step3")))

	req, _ := http.NewRequest(http.MethodGet, m.Server.URL+"/1", nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	assert.True(t, s1.IsDone())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step1", string(body))

	req, _ = http.NewRequest(http.MethodGet, m.Server.URL+"/2", nil)
	res, _ = http.DefaultClient.Do(req)
	body, _ = ioutil.ReadAll(res.Body)

	assert.True(t, s2.IsDone())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step2", string(body))

	req, _ = http.NewRequest(http.MethodGet, m.Server.URL+"/3", nil)
	res, _ = http.DefaultClient.Do(req)
	body, _ = ioutil.ReadAll(res.Body)

	assert.True(t, s3.IsDone())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step3", string(body))
}
