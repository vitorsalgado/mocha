package foundation_test

import "github.com/vitorsalgado/mocha/v3/foundation"

type testMock struct {
	*foundation.BaseMock
}

func newTestMock() *testMock {
	return &testMock{BaseMock: foundation.NewMock()}
}
