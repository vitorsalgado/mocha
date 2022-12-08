package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

type mockDebug struct {
	mock.Mock
}

func (m *mockDebug) OnError(err error) {
	m.Called(errors.Unwrap(err))
}

type mockRBP struct {
	fn func() error
}

func (m *mockRBP) CanParse(_ string, _ *http.Request) bool {
	return true
}

func (m *mockRBP) Parse(_ []byte, _ *http.Request) (any, error) {
	return nil, m.fn()
}

func TestMocha_Debug_SimpleError(t *testing.T) {
	expected := fmt.Errorf("boom")

	d := &mockDebug{}
	d.On("OnError", expected).Return(nil)

	rbp := &mockRBP{fn: func() error {
		return expected
	}}

	m := mocha.New(t,
		mocha.Configure().
			RequestBodyParsers(rbp).
			Debug(d.OnError))
	m.MustStart()

	defer m.Close()

	m.MustMock(mocha.Post(URLPath("/test")).
		ReplyJust(http.StatusOK, nil))

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader("hello"))
	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, mocha.StatusNoMockFound, res.StatusCode)
	d.AssertExpectations(t)
}

func TestMocha_Debug_Panic(t *testing.T) {
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	expected := fmt.Errorf("boom")

	rbp := &mockRBP{fn: func() error {
		panic(expected)
	}}

	spy := mock.Mock{}
	spy.On("func2", mock.Anything).Return(nil)

	fn := func(err error) {
		spy.Called(err)
		cancel()
	}

	m := mocha.New(t,
		mocha.Configure().
			RequestBodyParsers(rbp).
			Debug(fn))
	m.MustStart()

	defer m.Close()

	m.MustMock(mocha.Post(URLPath("/test")).
		ReplyJust(http.StatusOK, nil))

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader("hello"))
	res, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, mocha.StatusNoMockFound, res.StatusCode)

	<-ctx.Done()

	spy.AssertExpectations(t)
}
