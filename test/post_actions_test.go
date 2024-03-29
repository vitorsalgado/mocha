package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type action struct {
	mock.Mock
}

func (act *action) Run(a mocha.PostActionArgs) error {
	args := act.Called(a)
	return args.Error(0)
}

func TestPostAction(t *testing.T) {
	t.Parallel()

	t.Run("should call registered post action", func(t *testing.T) {
		m := mocha.New(t)
		m.Start()

		act := &action{}
		act.On("Run", mock.Anything).Return(nil)

		scope := m.AddMocks(mocha.Get(expect.URLPath("/test")).
			PostAction(act).
			Reply(reply.OK()))

		req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		scope.AssertCalled(t)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		act.AssertExpectations(t)
	})

	t.Run("should not be affected by errors on registered post actions", func(t *testing.T) {
		m := mocha.New(t)
		m.Start()

		defer m.Close()

		act := &action{}
		act.On("Run", mock.Anything).Return(fmt.Errorf("failed to run post action"))

		scope := m.AddMocks(mocha.Get(expect.URLPath("/test")).
			PostAction(act).
			Reply(reply.OK()))

		req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		scope.AssertCalled(t)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		act.AssertExpectations(t)
	})
}
