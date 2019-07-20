package jsl_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	jsl "github.com/json-schema-language/json-schema-language-go"
)

func Example() {
	schemaJSON := `
		{
			"properties": {
				"name": { "type": "string" },
				"age": { "type": "number" },
				"phones": {
					"elements": { "type": "string" }
				}
			}
		}
	`

	var schema jsl.Schema
	_ = json.Unmarshal([]byte(schemaJSON), &schema)

	inputOkJSON := `
		{
			"name": "John Doe",
			"age": 43,
			"phones": [
				"+44 1234567",
				"+44 2345678"
			]
		}
	`

	var inputOk interface{}
	_ = json.Unmarshal([]byte(inputOkJSON), &inputOk)

	inputBadJSON := `
		{
			"age": "43",
			"phones": [
					"+44 1234567",
					442345678
			]
		}
	`

	var inputBad interface{}
	_ = json.Unmarshal([]byte(inputBadJSON), &inputBad)

	validator := jsl.Validator{}
	resultOk, _ := validator.Validate(schema, inputOk)
	fmt.Println(resultOk.IsValid())

	resultBad, _ := validator.Validate(schema, inputBad)

	// To make tests predictable, we'll sort these errors in some arbitrary way.
	// You won't need to do this yourself in production, but you may need to do
	// something similar in your own tests.
	sort.Slice(resultBad.Errors, func(i, j int) bool {
		return strings.Join(resultBad.Errors[i].InstancePath, "/") < strings.Join(resultBad.Errors[j].InstancePath, "/")
	})

	fmt.Println(resultBad.Errors[0].InstancePath, resultBad.Errors[0].SchemaPath)
	fmt.Println(resultBad.Errors[1].InstancePath, resultBad.Errors[1].SchemaPath)
	fmt.Println(resultBad.Errors[2].InstancePath, resultBad.Errors[2].SchemaPath)

	// Output:
	// true
	// [] [properties name]
	// [age] [properties age type]
	// [phones 1] [properties phones elements type]
}
