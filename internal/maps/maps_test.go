package maps

import (
	"encoding/json"
	"github.com/vitorsalgado/mocha/internal/assert"
	"log"
	"testing"
)

func TestArray(t *testing.T) {
	jsonData := `
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
		},
		null
	]
]
`
	var data []any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal(err)
	}

	qa, err := Reach[string]("[0][1]", data)
	flg, err := Reach[float64]("[0][2].test", data)
	working, err := Reach[bool]("[0][2].entries[0][0].working", data)
	nilObj, err := Reach[any]("[0][3]", data)

	assert.Nil(t, err)
	assert.Equal(t, "qa", qa)
	assert.Equal(t, 1, flg)
	assert.True(t, working)
	assert.Nil(t, nilObj)
}

func TestObject(t *testing.T) {
	j := `
{
	"name": "test",
	"age": 100,
	"active": true,
	"jobs": ["qa", "dev", null],
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
	data := make(map[string]any)
	err := json.Unmarshal([]byte(j), &data)
	if err != nil {
		log.Fatal(err)
	}

	name, err := Reach[string]("name", data)
	age, err := Reach[float64]("age", data)
	active, err := Reach[bool]("active", data)
	ext, err := Reach[map[string]any]("extra", data)
	salary, err := Reach[float64]("extra.salary", data)
	employer, err := Reach[any]("extra.employer", data)
	street, err := Reach[string]("extra.address.street", data)
	jobs, err := Reach[[]any]("jobs", data)
	nilJob, err := Reach[any]("jobs[2]", data)
	qa, err := Reach[string]("jobs[0]", data)
	deepValue, err := Reach[string]("deep[1].params[0].name", data)
	nothing, err := Reach[any]("nothing", data)

	assert.Nil(t, err)
	assert.Equal(t, "test", name)
	assert.Equal(t, float64(100), age)
	assert.True(t, active)
	assert.NotNil(t, ext)
	assert.Equal(t, float64(50), ext["salary"].(float64))
	assert.Equal(t, 50, salary)
	assert.Nil(t, employer)
	assert.Equal(t, "somewhere nice", street)
	assert.Equal(t, 3, len(jobs))
	assert.Equal(t, []any{"qa", "dev", nil}, jobs)
	assert.Nil(t, nilJob)
	assert.Equal(t, "qa", qa)
	assert.Equal(t, "deep value", deepValue)
	assert.Nil(t, nothing)
}
