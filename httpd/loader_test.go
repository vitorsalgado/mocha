package mhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileLoader_Load(t *testing.T) {
	app := NewAPI(Setup().MockFilePatterns("testdata/loader/load/*mock.json", "testdata/loader/load/*.json"))
	loader := &fileLoader{}

	err := loader.Load(app)

	require.NoError(t, err)
	require.Equal(t, 2, len(app.storage.GetAll()))
}

func TestFileLoader_LoadWithError(t *testing.T) {
	app := NewAPI(Setup().MockFilePatterns("testdata/loader/invalid/*.json"))
	loader := &fileLoader{}

	err := loader.Load(app)

	require.Error(t, err)
}

func TestFileLoader(t *testing.T) {
	app := NewAPI(Setup().MockFilePatterns("testdata/loader/*.mock.*"))
	app.MustStart()

	defer app.Close()

	httpClient := &http.Client{}
	res, err := httpClient.Get(app.URL() + "/test?term=test&filter=all+none&page=10")
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Equal(t, "hi", string(b))

	res, err = httpClient.Get(app.URL() + "/test/002")
	require.NoError(t, err)

	b, err = io.ReadAll(res.Body)
	res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "test ok\n", string(b))

	res, err = httpClient.Get(app.URL() + "/test003?term=test&filter=none")
	require.NoError(t, err)

	body := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&body)
	res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Len(t, body, 0)
}
