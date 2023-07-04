package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzstd"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

func TestScenarios(t *testing.T) {
	testCases := []struct {
		name string
		s1   dzstd.Builder[*dzhttp.HTTPMock, *dzhttp.HTTPMockApp]
		s2   dzstd.Builder[*dzhttp.HTTPMock, *dzhttp.HTTPMockApp]
		s3   dzstd.Builder[*dzhttp.HTTPMock, *dzhttp.HTTPMockApp]
	}{
		{"code",
			dzhttp.Get(URLPath("/1")).
				StartScenario("code").
				ScenarioStateWillBe("step2").
				Name("step-1").
				Reply(dzhttp.OK().PlainText("step1")),
			dzhttp.Get(URLPath("/2")).
				ScenarioIs("code").
				ScenarioStateIs("step2").
				ScenarioStateWillBe("step3").
				Name("step-2").
				Reply(dzhttp.OK().PlainText("step2")),
			dzhttp.Get(URLPath("/3")).
				ScenarioIs("code").
				ScenarioStateIs("step3").
				ScenarioStateWillBe("step4").
				Name("step-3").
				Reply(dzhttp.OK().PlainText("step3"))},

		{"file",
			dzhttp.FromFile("testdata/scenario/step_1.yaml"),
			dzhttp.FromFile("testdata/scenario/step_2.yaml"),
			dzhttp.FromFile("testdata/scenario/step_3.yaml")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := dzhttp.NewAPIWithT(t)
			m.MustStart()

			s1 := m.MustMock(tc.s1)
			s2 := m.MustMock(tc.s2)
			s3 := m.MustMock(tc.s3)

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
			require.Equal(t, dzhttp.StatusNoMatch, res.StatusCode)

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
			require.Equal(t, dzhttp.StatusNoMatch, res.StatusCode)

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
			require.Equal(t, dzhttp.StatusNoMatch, res.StatusCode)
		})
	}
}
