package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestScenarios(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	scn := "test"

	s1 := m.MustMock(mocha.Get(matcher.URLPath("/1")).
		StartScenario(scn).
		ScenarioStateWillBe("step2").
		Name("step-1").
		Reply(mocha.OK().PlainText("step1")))

	s2 := m.MustMock(mocha.Get(matcher.URLPath("/2")).
		ScenarioIs(scn).
		ScenarioStateIs("step2").
		ScenarioStateWillBe("step3").
		Name("step-2").
		Reply(mocha.OK().PlainText("step2")))

	s3 := m.MustMock(mocha.Get(matcher.URLPath("/3")).
		ScenarioIs(scn).
		ScenarioStateIs("step3").
		ScenarioStateWillBe("step4").
		Name("step-3").
		Reply(mocha.OK().PlainText("step3")))

	// --- step1

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/1", nil)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.True(t, s1.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "step1", string(body))
	require.NoError(t, res.Body.Close())

	// step1 is already in a different state, it should not match anymore.
	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/1", nil)
	res, err = http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)

	// --- step2

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/2", nil)
	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	require.True(t, s2.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "step2", string(body))
	require.NoError(t, res.Body.Close())

	// step2 is already in a different state, it should not match anymore.
	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/2", nil)
	res, err = http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)

	// --- step3

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/3", nil)
	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	require.True(t, s3.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "step3", string(body))
	require.NoError(t, res.Body.Close())

	// step3 is already in a different state, it should not match anymore.
	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/3", nil)
	res, err = http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}

func TestScenarios_SetupFromFile(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	s1 := m.MustMock(mocha.FromFile("testdata/scenario/step_1.yaml"))
	s2 := m.MustMock(mocha.FromFile("testdata/scenario/step_2.yaml"))
	s3 := m.MustMock(mocha.FromFile("testdata/scenario/step_3.yaml"))

	// --- step1

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/1", nil)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.True(t, s1.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "step1", string(body))
	require.NoError(t, res.Body.Close())

	// step1 is already in a different state, it should not match anymore.
	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/1", nil)
	res, err = http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)

	// --- step2

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/2", nil)
	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	require.True(t, s2.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "step2", string(body))
	require.NoError(t, res.Body.Close())

	// step2 is already in a different state, it should not match anymore.
	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/2", nil)
	res, err = http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)

	// --- step3

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/3", nil)
	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	require.True(t, s3.HasBeenCalled())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "step3", string(body))
	require.NoError(t, res.Body.Close())

	// step3 is already in a different state, it should not match anymore.
	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/3", nil)
	res, err = http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}
