package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type action struct {
	mock.Mock
}

func (act *action) Run(a *mocha.PostActionInput) error {
	args := act.Called(a)
	return args.Error(0)
}

func TestPostAction(t *testing.T) {
	t.Run("should call registered post action", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		act := &action{}
		act.On("Run", mock.Anything).Return(nil)

		scope := m.MustMock(mocha.Get(matcher.URLPath("/test")).
			PostAction(act).
			Reply(mocha.OK()))

		req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
		res, err := req.Do()

		assert.NoError(t, err)
		scope.AssertCalled(t)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		act.AssertExpectations(t)
	})

	t.Run("should not be affected by errors on registered post actions", func(t *testing.T) {
		m := mocha.New()
		m.MustStart()

		defer m.Close()

		act := &action{}
		act.On("Run", mock.Anything).Return(fmt.Errorf("failed to run post action"))

		scope := m.MustMock(mocha.Get(matcher.URLPath("/test")).
			PostAction(act).
			Reply(mocha.OK()))

		req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
		res, err := req.Do()

		assert.NoError(t, err)
		scope.AssertCalled(t)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		act.AssertExpectations(t)
	})
}
