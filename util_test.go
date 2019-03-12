package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJSONFile(t *testing.T) {
	assert := assert.New(t)

	type Person struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Friends []Person `json:"friends"`
	}

	person := Person{}

	if err := ParseJSONFile("./testfiles/person.json", &person); err != nil {
		assert.Fail(err.Error())
		return
	}

	t.Log(person)
	assert.Equal("Yunhong", person.Name)
	assert.Equal(12, person.Age)
	assert.ElementsMatch([]Person{
		Person{Name: "drama1", Age: 20},
		Person{Name: "drama2", Age: 25},
	}, person.Friends)
}
