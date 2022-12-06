package mocha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoader_Load(t *testing.T) {
	app := New(t, Configure().MockFilePatterns("testdata/0/*mock.json", "testdata/0/*.json").Build())
	loader := &FileLoader{}

	err := loader.Load(app)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(app.storage.FetchAll()))
}
