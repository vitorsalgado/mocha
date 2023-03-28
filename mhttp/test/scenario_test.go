package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/foundation"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
)

func TestScenarios(t *testing.T) {
	testCases := []struct {
		name string
		s1   foundation.Builder[*mhttp2.HTTPMock, *mhttp2.HTTPMockApp]
		s2   foundation.Builder[*mhttp2.HTTPMock, *mhttp2.HTTPMockApp]
		s3   foundation.Builder[*mhttp2.HTTPMock, *mhttp2.HTTPMockApp]
	}{
		{"code",
			mhttp2.Get(URLPath("/1")).
				StartScenario("code").
				ScenarioStateWillBe("step2").
				Name("step-1").
				Reply(mhttp2.OK().PlainText("step1")),
			mhttp2.Get(URLPath("/2")).
				ScenarioIs("code").
				ScenarioStateIs("step2").
				ScenarioStateWillBe("step3").
				Name("step-2").
				Reply(mhttp2.OK().PlainText("step2")),
			mhttp2.Get(URLPath("/3")).
				ScenarioIs("code").
				ScenarioStateIs("step3").
				ScenarioStateWillBe("step4").
				Name("step-3").
				Reply(mhttp2.OK().PlainText("step3"))},

		{"file",
			mhttp2.FromFile("testdata/scenario/step_1.yaml"),
			mhttp2.FromFile("testdata/scenario/step_2.yaml"),
			mhttp2.FromFile("testdata/scenario/step_3.yaml")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mhttp2.NewAPIWithT(t)
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
			require.Equal(t, mhttp2.StatusNoMatch, res.StatusCode)

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
			require.Equal(t, mhttp2.StatusNoMatch, res.StatusCode)

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
			require.Equal(t, mhttp2.StatusNoMatch, res.StatusCode)
		})
	}
}
