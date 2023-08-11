package dzstd

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestInterceptor struct{ txt string }

func (in *TestInterceptor) Intercept(p []byte, chain Chain) (n int, err error) {
	return chain.Next(append(p, []byte("--"+in.txt)...))
}

type TestW struct {
	mock.Mock
}

func (m *TestW) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestInterceptors(t *testing.T) {
	data := []byte("init")
	expected := []byte("init--hello--world")
	in1 := &TestInterceptor{"hello"}
	in2 := &TestInterceptor{"world"}

	w := &TestW{}
	w.On("Write", expected).Return(len(expected), nil)
	root := &RootIntereptor{w}
	chain := &InterceptorChain{interceptors: []Interceptor{in1, in2, root}}

	n, err := chain.Next(data)

	require.NoError(t, err)
	require.Equal(t, len(expected), n)
	require.True(t, w.AssertCalled(t, "Write", expected))
}
