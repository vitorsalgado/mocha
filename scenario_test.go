package mocha

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/reply"
)

func TestScenario(t *testing.T) {
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
