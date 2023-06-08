package lib_test

import "github.com/vitorsalgado/mocha/v3/lib"

type testMock struct {
	*lib.BaseMock
}

func newTestMock() *testMock {
	return &testMock{BaseMock: lib.NewMock()}
}
