package jsonpath

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
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

	c, err := Get("test", jsonData)
	assert.Nil(t, c)
	assert.Error(t, err)

	i, err := Get("[3][1]", jsonData)
	assert.Nil(t, i)
	assert.Error(t, err)

	qa, err := Get("[0][1]", data)
	assert.Nil(t, err)
	assert.Equal(t, "qa", qa)

	flg, err := Get("[0][2].test", data)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), flg.(float64))

	working, err := Get("[0][2].entries[0][0].working", data)
	assert.Nil(t, err)
	assert.True(t, working.(bool))

	i2, err := Get("[0][2].entries[0][2].working", data)
	assert.Nil(t, i2)
	assert.NotNil(t, err)

	nilObj, err := Get("[0][3]", data)
	assert.Nil(t, err)
	assert.Nil(t, nilObj)

	noIdx, err := Get("[-1]", data)
	assert.Nil(t, noIdx)
	assert.NotNil(t, err)

	nan, err := Get("[abc]", data)
	assert.Nil(t, nan)
	assert.NotNil(t, err)
}

func TestObject(t *testing.T) {
	jsonData := `
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
		},
		null
	],
	"nothing": null
}
`
	mapped := make(map[string]any)
	err := json.Unmarshal([]byte(jsonData), &mapped)
	if err != nil {
		log.Fatal(err)
	}

	br, err := Get("[1].name", jsonData)
	assert.Nil(t, br)
	assert.Error(t, err)

	name, err := Get("name", mapped)
	assert.Nil(t, err)
	assert.Equal(t, "test", name)

	age, err := Get("age", mapped)
	assert.Nil(t, err)
	assert.Equal(t, float64(100), age.(float64))

	active, err := Get("active", mapped)
	assert.Nil(t, err)
	assert.True(t, active.(bool))

	ext, err := Get("extra", mapped)
	assert.Nil(t, err)
	assert.NotNil(t, ext)
	assert.Equal(t, float64(50), ext.(map[string]any)["salary"].(float64))

	salary, err := Get("extra.salary", mapped)
	assert.Nil(t, err)

	assert.Equal(t, float64(50), salary.(float64))

	employer, err := Get("extra.employer", mapped)
	assert.Nil(t, err)
	assert.Nil(t, employer)

	street, err := Get("extra.address.street", mapped)
	assert.Nil(t, err)
	assert.Equal(t, "somewhere nice", street)

	jobs, err := Get("jobs", mapped)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(jobs.([]any)))
	assert.Equal(t, []any{"qa", "dev", nil}, jobs.([]any))

	nilJob, err := Get("jobs[2]", mapped)
	assert.Nil(t, err)
	assert.Nil(t, nilJob)

	qa, err := Get("jobs[0]", mapped)
	assert.Nil(t, err)
	assert.Equal(t, "qa", qa)

	deepValue, err := Get("deep[1].params[0].name", mapped)
	assert.Nil(t, err)
	assert.Equal(t, "deep value", deepValue)

	nullDeepValue, err := Get("deep[2]", mapped)
	assert.Nil(t, err)
	assert.Nil(t, nullDeepValue)

	nonExistentDeepValue, err := Get("deep[3]", mapped)
	assert.NotNil(t, err)
	assert.Nil(t, nonExistentDeepValue)

	nothing, err := Get("nothing", mapped)
	assert.Nil(t, err)
	assert.Nil(t, nothing)

	// not found
	no, err := Get("not_present", mapped)
	assert.Nil(t, no)
	assert.NotNil(t, err)
}
