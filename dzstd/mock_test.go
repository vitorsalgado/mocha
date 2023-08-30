package dzstd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemoryStorage(t *testing.T) {
	st := NewStore[*BaseMock]()
	ctx := context.Background()

	mock1 := &BaseMock{ID: "1", Name: "mock_1", enabled: true, Priority: 0}
	mock2 := &BaseMock{ID: "2", Name: "mock_2", enabled: true, Priority: 1}
	mock3 := &BaseMock{ID: "3", Name: "mock_3", enabled: true, Priority: 2}

	err := st.Save(ctx, mock1)
	require.NoError(t, err)

	err = st.Save(ctx, mock2)
	require.NoError(t, err)

	err = st.Save(ctx, mock3)
	require.NoError(t, err)

	mock, err := st.Get(ctx, mock1.ID)
	require.NoError(t, err)
	require.Equal(t, mock1, mock)

	mocks, err := st.GetAll(ctx)
	require.NoError(t, err)

	m := mocks[0]
	require.Equal(t, m.ID, "1")
	require.Equal(t, m.Name, "mock_1")

	m.Disable()

	eligible := make([]*BaseMock, 0)
	done := make(chan struct{})
	out, err := st.FindEligible(ctx, done)
	for v := range out {
		eligible = append(eligible, v)
	}

	require.NoError(t, err)
	require.Len(t, eligible, 2)

	eligible, err = st.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, eligible, 3)

	err = st.Delete(ctx, "2")
	require.NoError(t, err)
	eligible, err = st.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, eligible, 2)

	eligible = make([]*BaseMock, 0)
	out, err = st.FindEligible(ctx, done)
	for v := range out {
		eligible = append(eligible, v)
	}

	require.NoError(t, err)
	require.Len(t, eligible, 1)
	require.Equal(t, "3", eligible[0].ID)

	err = st.DeleteAll(ctx)
	require.NoError(t, err)

	eligible, err = st.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, eligible, 0)

	err = st.Save(ctx, &BaseMock{ID: "10", Name: "mock_ext_1", enabled: true, Priority: 0, Source: "ext"})
	require.NoError(t, err)
	err = st.Save(ctx, &BaseMock{ID: "11", Name: "mock_11", enabled: true, Priority: 0})
	require.NoError(t, err)

	all, err := st.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, all, 2)

	err = st.DeleteExternal(ctx)
	require.NoError(t, err)

	eligible, err = st.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, eligible, 1)
	require.Equal(t, "11", eligible[0].ID)

	err = st.DeleteAll(ctx)
	require.NoError(t, err)

	all, err = st.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, all, 0)

	err = st.DeleteExternal(ctx)
	require.NoError(t, err)
}

func TestDeleteExt(t *testing.T) {
	ctx := context.Background()
	st := NewStore[*BaseMock]()

	require.NoError(t, st.Save(ctx, &BaseMock{ID: "10", Name: "mock_ext_1", enabled: true, Priority: 0, Source: "ext"}))
	require.NoError(t, st.Save(ctx, &BaseMock{ID: "10", Name: "mock_ext_2", enabled: true, Priority: 0, Source: "ext"}))
	require.NoError(t, st.DeleteExternal(ctx))

	all, err := st.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, all, 0)
}
