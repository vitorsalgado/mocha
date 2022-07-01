package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/matcher"
	mok "github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/reply"
)

type action struct {
	mock.Mock
}

func (act *action) Run(a mok.PostActionArgs) error {
	args := act.Called(a)
	return args.Error(0)
}

func TestPostAction(t *testing.T) {
	t.Parallel()

	t.Run("should call registered post action", func(t *testing.T) {
		m := mocha.ForTest(t)
		m.Start()

		act := &action{}
		act.On("Run", mock.Anything).Return(nil)

		scope := m.Mock(mocha.Get(matcher.URLPath("/test")).
			PostAction(act).
			Reply(reply.OK()))

		req := testutil.Get(fmt.Sprintf("%s/test", m.Server.URL))
		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, scope.Done())
		assert.Equal(t, http.StatusOK, res.StatusCode)
		act.AssertExpectations(t)
	})

	t.Run("should not be affected by errors on registered post actions", func(t *testing.T) {

		m := mocha.ForTest(t)
		m.Start()

		act := &action{}
		act.On("Run", mock.Anything).Return(fmt.Errorf("failed to run post action"))

		scope := m.Mock(mocha.Get(matcher.URLPath("/test")).
			PostAction(act).
			Reply(reply.OK()))

		req := testutil.Get(fmt.Sprintf("%s/test", m.Server.URL))
		res, err := req.Do()
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, scope.Done())
		assert.Equal(t, http.StatusOK, res.StatusCode)
		act.AssertExpectations(t)
	})
}
