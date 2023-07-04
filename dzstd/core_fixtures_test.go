package dzstd_test

import "github.com/vitorsalgado/mocha/v3/dzstd"

type testMock struct {
	*dzstd.BaseMock
}

func newTestMock() *testMock {
	return &testMock{BaseMock: dzstd.NewMock()}
}
