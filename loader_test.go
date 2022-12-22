package mocha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoader_Load(t *testing.T) {
	app := New(
		Configure().Dirs("testdata/0/*mock.json", "testdata/0/*.json"))
	loader := &FileLoader{}

	err := loader.Load(app)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(app.storage.GetAll()))
}
