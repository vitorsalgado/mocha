package mocha

import (
	"github.com/vitorsalgado/mocha/internal/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestScenario(t *testing.T) {
	m := NewT(t)

	s1 := m.Mock(Get(URLPath("/1")).
		Scenario("test", ScenarioStarted, "step2").
		Name("step-1").
		Reply(OK().BodyStr("step1")))

	s2 := m.Mock(Get(URLPath("/2")).
		Scenario("test", "step2", "step3").
		Name("step-2").
		Reply(OK().BodyStr("step2")))

	s3 := m.Mock(Get(URLPath("/3")).
		Scenario("test", "step3", "step4").
		Name("step-3").
		Reply(OK().BodyStr("step3")))

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
