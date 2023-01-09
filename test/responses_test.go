package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
)

func TestJSONResponse(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	type payload struct {
		Language string `json:"language"`
		Active   bool   `json:"active"`
	}

	p := &payload{"go", true}

	m.MustMock(mocha.Getf("/test").Reply(mocha.OK().JSON(p)))

	res, err := testutil.Get(m.URL() + "/test").Do()
	require.NoError(t, err)

	defer res.Body.Close()

	var body payload

	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, p, &body)
}
