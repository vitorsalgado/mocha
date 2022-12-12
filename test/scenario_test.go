package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestScenarioMatcher(t *testing.T) {
	m := mocha.New(t)
	m.MustStart()

	defer m.Close()

	scn := "test"

	s1 := m.MustMock(mocha.Get(matcher.URLPath("/1")).
		StartScenario(scn).
		ScenarioStateWillBe("step2").
		Name("step-1").
		Reply(reply.OK().PlainText("step1")))

	s2 := m.MustMock(mocha.Get(matcher.URLPath("/2")).
		ScenarioIs(scn).
		ScenarioStateIs("step2").
		ScenarioStateWillBe("step3").
		Name("step-2").
		Reply(reply.OK().PlainText("step2")))

	s3 := m.MustMock(mocha.Get(matcher.URLPath("/3")).
		ScenarioIs(scn).
		ScenarioStateIs("step3").
		ScenarioStateWillBe("step4").
		Name("step-3").
		Reply(reply.OK().PlainText("step3")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/1", nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	assert.True(t, s1.HasBeenCalled())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step1", string(body))

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/2", nil)
	res, _ = http.DefaultClient.Do(req)
	body, _ = io.ReadAll(res.Body)

	assert.True(t, s2.HasBeenCalled())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step2", string(body))

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/3", nil)
	res, _ = http.DefaultClient.Do(req)
	body, _ = io.ReadAll(res.Body)

	assert.True(t, s3.HasBeenCalled())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "step3", string(body))
}
