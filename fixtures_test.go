package mocha

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

var _req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

func newReqValues(req *http.Request) *RequestValues {
	return &RequestValues{RawRequest: req, URL: req.URL}
}

var _ TestingT = (*fakeT)(nil)

type fakeT struct{ mock.Mock }

func newFakeT() *fakeT {
	t := &fakeT{}
	t.On("Helper").Return()
	t.On("Logf", mock.Anything, mock.Anything).Return()
	t.On("Errorf", mock.Anything, mock.Anything).Return()
	t.On("FailNow").Return()

	return t
}

func (m *fakeT) Helper() {
	m.Called()
}

func (m *fakeT) Logf(format string, args ...any) {
	m.Called(format, args)
}

func (m *fakeT) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func (m *fakeT) Cleanup(_ func()) {}
