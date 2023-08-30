package dzhttp

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestRace(t *testing.T) {
	m := newMock()
	jobs := 10
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func() {
			m.Inc()
			m.Hits()
			wg.Done()
		}()

		m.Inc()
	}

	m.Hits()
	m.Inc()
	m.Inc()

	wg.Wait()

	assert.EqualValues(t, (jobs*2)+2, m.Hits())
}

func TestMock(t *testing.T) {
	m := newMock()

	t.Run("should init enabled", func(t *testing.T) {
		assert.True(t, m.IsEnabled())
	})

	t.Run("should disable mock when calling .Disable()", func(t *testing.T) {
		m.Disable()
		assert.False(t, m.IsEnabled())

		m.Enable()
		assert.True(t, m.IsEnabled())
	})

	t.Run("should return called when it was hit", func(t *testing.T) {
		assert.False(t, m.HasBeenCalled())
		m.Inc()
		assert.True(t, m.HasBeenCalled())

		m.Dec()
		assert.False(t, m.HasBeenCalled())
	})
}

func TestMockBuild(t *testing.T) {
	m := newMock()
	m.Inc()
	m.Disable()

	mm, err := m.Build()

	assert.NoError(t, err)
	assert.Equal(t, m, mm)
}

func TestMock_MarshalJSON(t *testing.T) {
	app := NewAPI()
	builder := Request().
		Header("accept", Eq("test")).
		Once().
		ScenarioIs("test").ScenarioStateIs("state").ScenarioStateWillBe("future").
		Reply(OK())

	m, err := builder.Build(app)
	require.NoError(t, err)

	b, err := m.MarshalJSON()
	require.NoError(t, err)

	fmt.Println(string(b))
}
