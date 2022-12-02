package mocha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoader_Load(t *testing.T) {
	app := New(t, Configure().Pattern("testdata/d/*mock.json").Build())
	loader := &FileLoader{}

	err := loader.Load(app)

	assert.NoError(t, err)
	assert.Len(t, app.storage.FetchAll(), 1)
}
