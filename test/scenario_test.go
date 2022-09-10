package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v2"
	"github.com/vitorsalgado/mocha/v2/expect"
	"github.com/vitorsalgado/mocha/v2/reply"
)

func TestScenarioMatcher(t *testing.T) {
	m := mocha.New(t)
	m.Start()
	scn := "test"

	s1 := m.AddMocks(mocha.Get(expect.URLPath("/1")).
		StartScenario(scn).
		ScenarioStateWillBe("step2").
		Name("step-1").
		Reply(reply.OK().BodyString("step1")))

	s2 := m.AddMocks(mocha.Get(expect.URLPath("/2")).
		ScenarioIs(scn).
		ScenarioStateIs("step2").
		ScenarioStateWillBe("step3").
		Name("step-2").
		Reply(reply.OK().BodyString("step2")))

	s3 := m.AddMocks(mocha.Get(expect.URLPath("/3")).
		ScenarioIs(scn).
		ScenarioStateIs("step3").
		ScenarioStateWillBe("step4").
		Name("step-3").
		Reply(reply.OK().BodyString("step3")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/1", nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	assert.True(t, s1.Called())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step1", string(body))

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/2", nil)
	res, _ = http.DefaultClient.Do(req)
	body, _ = io.ReadAll(res.Body)

	assert.True(t, s2.Called())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step2", string(body))

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/3", nil)
	res, _ = http.DefaultClient.Do(req)
	body, _ = io.ReadAll(res.Body)

	assert.True(t, s3.Called())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step3", string(body))
}
