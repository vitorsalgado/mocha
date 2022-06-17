package maps

import (
	"encoding/json"
	"github.com/vitorsalgado/mocha/internal/assert"
	"log"
	"testing"
)

func TestGetNew(t *testing.T) {
	j := `
[
	[
		"dev", 
		"qa", 
		{
			"test": 1,
			"entries": [
				[
					{ "working": true }
				]
			]
		}
	]
]
`
	var data []any
	err := json.Unmarshal([]byte(j), &data)
	if err != nil {
		log.Fatal(err)
	}

	qa, err := GetNew[string]("[0][1]", data)
	working, err := GetNew[bool]("[0][2].entries[0][0].working", data)

	assert.Nil(t, err)
	assert.Equal(t, "qa", qa)
	//assert.Equal(t, 1, testOk)
	assert.True(t, working)
}

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

	name, _ := GetNew[string]("name", m)
	age, _ := GetNew[float64]("age", m)
	active, _ := GetNew[bool]("active", m)
	salary, _ := GetNew[float64]("extra.salary", m)
	employer, _ := GetNew[any]("extra.employer", m)
	street, _ := GetNew[string]("extra.address.street", m)
	jobs, _ := GetNew[[]any]("jobs", m)
	qa, _ := GetNew[string]("jobs[0]", m)
	deepValue, _ := GetNew[string]("deep[1].params[0].name", m)
	nothing, _ := GetNew[any]("nothing", m)

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
