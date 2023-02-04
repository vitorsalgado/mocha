package mocha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryStorage(t *testing.T) {
	st := newStore()

	mock1 := &Mock{ID: "1", Name: "mock_1", Enabled: true, Priority: 0}
	mock2 := &Mock{ID: "2", Name: "mock_2", Enabled: true, Priority: 1}
	mock3 := &Mock{ID: "3", Name: "mock_3", Enabled: true, Priority: 2}

	st.Save(mock1)
	st.Save(mock2)
	st.Save(mock3)

	assert.Equal(t, mock1, st.Get(mock1.ID))

	m := st.GetAll()[0]
	assert.Equal(t, m.ID, "1")
	assert.Equal(t, m.Name, "mock_1")

	m.Disable()

	mocks := st.GetEligible()
	assert.Len(t, mocks, 2)

	mocks = st.GetAll()
	assert.Len(t, mocks, 3)

	st.Delete("2")
	mocks = st.GetAll()
	assert.Len(t, mocks, 2)

	mocks = st.GetEligible()
	assert.Len(t, mocks, 1)
	assert.Equal(t, "3", mocks[0].ID)

	st.DeleteAll()

	mocks = st.GetAll()

	assert.Len(t, mocks, 0)

	st.Save(&Mock{ID: "10", Name: "mock_ext_1", Enabled: true, Priority: 0, Source: "ext"})
	st.Save(&Mock{ID: "11", Name: "mock_11", Enabled: true, Priority: 0})

	assert.Len(t, st.GetAll(), 2)

	st.DeleteExternal()

	mocks = st.GetAll()

	assert.Len(t, mocks, 1)
	assert.Equal(t, "11", mocks[0].ID)
}
