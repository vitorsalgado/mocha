package mocha

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

func TestBuilderFs_CustomMockFileHandlers(t *testing.T) {
	m := NewAPI(Setup().MockFileHandlers(&customMockFileHandler{}))
	m.MustStart()

	defer m.Close()

	m.MustMock(FromFile("testdata/builder_fs/mock_file_handler/fixture.json"))

	httpClient := &http.Client{}
	res, err := httpClient.Get(m.URL() + "/test?term=test&filter=all&page=10&q=hello-world&custom=100")
	require.NoError(t, err)

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Equal(t, "hi", string(b))
}

func TestBuilderFs_FromBytes(t *testing.T) {
	testCases := []struct {
		filename  string
		extension string
	}{
		{"testdata/builder_fs/from_bytes/fixture_1.yaml", "yaml"},
		{"testdata/builder_fs/from_bytes/fixture_1.json", "json"},
	}

	httpClient := &http.Client{}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			m := NewAPI()
			m.MustStart()

			file, err := os.Open(tc.filename)
			require.NoError(t, err)

			b, err := io.ReadAll(file)
			require.NoError(t, err)

			m.MustMock(FromBytes(b, tc.extension))

			res, err := httpClient.Get(m.URL() + "/test")

			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)

			m.Close()
			file.Close()
		})
	}
}

func TestBuilderFs_MustNotAllowMultipleReplyDefinitions(t *testing.T) {
	m := NewAPI()

	filenames := []string{
		"testdata/builder_fs/multi_reply_def/1_multi_reply.yaml",
		"testdata/builder_fs/multi_reply_def/2_multi_reply.yaml",
		"testdata/builder_fs/multi_reply_def/3_multi_reply.yaml",
	}

	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			s, err := m.Mock(FromFile(filename))

			require.Nil(t, s)
			require.Error(t, err)
		})
	}
}

func TestBuilderFs_InvalidSchemaValid(t *testing.T) {
	m := NewAPI()
	filenames := []string{
		"testdata/builder_fs/invalid_schema/1_invalid_schema.yaml",
		"testdata/builder_fs/invalid_schema/2_no_request.yaml",
	}

	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			s, err := m.Mock(FromFile(filename))

			require.Nil(t, s)
			require.Error(t, err)
		})
	}
}
