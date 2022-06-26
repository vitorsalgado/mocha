package mock

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestInMemoryStorage(t *testing.T) {
	storage := NewStorage()
	storage.Save(&Mock{ID: 1, Name: "mock_1", Enabled: true, mu: &sync.Mutex{}})
	storage.Save(&Mock{ID: 2, Name: "mock_2", Enabled: true, mu: &sync.Mutex{}})
	storage.Save(&Mock{ID: 3, Name: "mock_3", Enabled: true, mu: &sync.Mutex{}})

	m := storage.FetchByID(1)
	assert.Equal(t, m.ID, 1)
	assert.Equal(t, m.Name, "mock_1")

	mocks := storage.FetchByIDs(1, 2)
	assert.Len(t, mocks, 2)

	m.Disable()

	mocks = storage.FetchEligibleSorted()
	assert.Len(t, mocks, 2)

	mocks = storage.FetchAll()
	assert.Len(t, mocks, 3)

	storage.Delete(2)
	mocks = storage.FetchAll()
	assert.Len(t, mocks, 2)

	mocks = storage.FetchEligibleSorted()
	assert.Len(t, mocks, 1)
	assert.Equal(t, 3, mocks[0].ID)

	storage.Flush()

	mocks = storage.FetchAll()

	assert.Len(t, mocks, 0)
}
