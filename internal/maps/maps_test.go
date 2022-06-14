package maps

import (
	"encoding/json"
	"github.com/vitorsalgado/mocha/internal/assert"
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	j := `
{
	"name": "test",
	"age": 100,
	"active": true,
	"jobs": ["qa", "dev"],
	"extra": {
		"salary": 50,
		"home": "Chile",
		"employer": null,
		"address": {
			"street": "somewhere nice"
		}
	},
	"deep": [
		{},
		{
			"key": "001",
			"params": [{ "name": "deep value" }]
		}
	],
	"nothing": null
}
`
	m := make(map[string]any)
	err := json.Unmarshal([]byte(j), &m)
	if err != nil {
		log.Fatal(err)
	}

	name := Get[string]("name", m)
	age := Get[float64]("age", m)
	active := Get[bool]("active", m)
	salary := Get[float64]("extra.salary", m)
	employer := Get[any]("extra.employer", m)
	street := Get[string]("extra.address.street", m)
	jobs := Get[[]any]("jobs", m)
	qa := Get[string]("jobs[0]", m)
	deepValue := Get[string]("deep[1].params[0].name", m)
	nothing := Get[any]("nothing", m)

	assert.Equal(t, "test", name)
	assert.Equal(t, float64(100), age)
	assert.True(t, active)
	assert.Equal(t, 50, salary)
	assert.Nil(t, employer)
	assert.Equal(t, "somewhere nice", street)
	assert.Equal(t, 2, len(jobs))
	assert.Equal(t, []any{"qa", "dev"}, jobs)
	assert.Equal(t, "qa", qa)
	assert.Equal(t, "deep value", deepValue)
	assert.Nil(t, nothing)
}
