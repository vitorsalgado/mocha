package jsonpath

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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

	qa, err := Get("[0][1]", data)
	assert.Nil(t, err)
	assert.Equal(t, "qa", qa)

	flg, err := Get("[0][2].test", data)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), flg.(float64))

	working, err := Get("[0][2].entries[0][0].working", data)
	assert.Nil(t, err)
	assert.True(t, working.(bool))

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
		},
		null
	],
	"nothing": null
}
`
	data := make(map[string]any)
	err := json.Unmarshal([]byte(j), &data)
	if err != nil {
		log.Fatal(err)
	}

	name, err := Get("name", data)
	assert.Nil(t, err)
	age, err := Get("age", data)
	assert.Nil(t, err)
	active, err := Get("active", data)
	assert.Nil(t, err)
	ext, err := Get("extra", data)
	assert.Nil(t, err)
	salary, err := Get("extra.salary", data)
	assert.Nil(t, err)
	employer, err := Get("extra.employer", data)
	assert.Nil(t, err)
	street, err := Get("extra.address.street", data)
	assert.Nil(t, err)
	jobs, err := Get("jobs", data)
	assert.Nil(t, err)
	nilJob, err := Get("jobs[2]", data)
	assert.Nil(t, err)
	qa, err := Get("jobs[0]", data)
	assert.Nil(t, err)
	deepValue, err := Get("deep[1].params[0].name", data)
	assert.Nil(t, err)
	nullDeepValue, err := Get("deep[2]", data)
	assert.Nil(t, err)
	nonExistentDeepValue, err := Get("deep[3]", data)
	assert.NotNil(t, err)
	nothing, err := Get("nothing", data)
	assert.Nil(t, err)

	// not found
	no, err := Get("not_present", data)
	assert.Nil(t, no)
	assert.NotNil(t, err)

	assert.Equal(t, "test", name)
	assert.Equal(t, float64(100), age.(float64))
	assert.True(t, active.(bool))
	assert.NotNil(t, ext)
	assert.Equal(t, float64(50), ext.(map[string]any)["salary"].(float64))
	assert.Equal(t, float64(50), salary.(float64))
	assert.Nil(t, employer)
	assert.Equal(t, "somewhere nice", street)
	assert.Equal(t, 3, len(jobs.([]any)))
	assert.Equal(t, []any{"qa", "dev", nil}, jobs.([]any))
	assert.Nil(t, nilJob)
	assert.Equal(t, "qa", qa)
	assert.Equal(t, "deep value", deepValue)
	assert.Nil(t, nullDeepValue)
	assert.Nil(t, nonExistentDeepValue)
	assert.Nil(t, nothing)
}
