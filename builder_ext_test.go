package mocha

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type customMockFileHandler struct {
}

func (p *customMockFileHandler) Handle(v map[string]any, b *MockBuilder) error {
	custom := v["custom"]
	query := custom.(map[string]any)["query"]
	key := query.(map[string]any)["key"]
	value := query.(map[string]any)["value"]

	b.Queryf(key.(string), value.(string))

	another := v["another_custom_field"]

	b.Queryf("custom", fmt.Sprintf("%v", another))

	return nil
}

func TestCustomMockFileHandlers(t *testing.T) {
	m := New(Configure().MockFileHandlers(&customMockFileHandler{}))
	m.MustStart()

	defer m.Close()

	m.MustMock(MockFromFile("testdata/mock_file_handler/handler.json"))

	httpClient := &http.Client{}
	res, err := httpClient.Get(m.URL() + "/test?term=test&filter=all&page=10&q=hello-world&custom=100")
	require.NoError(t, err)

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Equal(t, "hi", string(b))
}
