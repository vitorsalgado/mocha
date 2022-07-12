package mocha

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
)

func TestExpectations(t *testing.T) {
	m := New(t)
	m.Start()

	scoped := m.Mock(
		Get(expect.URLPath("/test")).
			Cond(Expect(Header("hello")).ToEqual("world")).
			Reply(reply.OK()))

	req := testutil.Get(fmt.Sprintf("%s/test", m.URL()))
	req.Header("hello", "world")
	res, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}

	scoped.AssertCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
