# jsl [![badge]][docs]

> Documentation on godoc.org: https://godoc.org/github.com/json-schema-language/json-schema-language-go

[badge]: https://godoc.org/github.com/json-schema-language/json-schema-language-go?status.svg
[docs]: https://godoc.org/github.com/json-schema-language/json-schema-language-go

This package is a Golang implementation of **JSON Schema Langugage**. You can
use this package to:

* Validate input data against a schema,
* Get a list of validation errors from that input data, or
* Build your own tooling on top of JSON Schema Language

## About JSON Schema Language

JSON Schema Language ("JSL") is a JSON-based format for describing the shape of
JSON data. You define the shape of your data using schemas, and from those
schemas you can:

* Validate that inputted data is correct against the schema
* Get and transmit standardized error details when data fails to validate
* Generate code, documentation, and user interfaces from schemas

This package fits into the picture by providing Golang support for parsing
schemas and checking data against schemas. Check out the [json-schema-language
website][jsl-website] for other tooling built on JSON Schema Language.

## Usage

See [the docs] for more detailed usage, but at a high level, here's how you
parse schemas and validate input data against them:

```golang
import (
  "encoding/json"
  "fmt"

  // It's recommended that you import this package as "jsl".
  jsl "github.com/json-schema-language/json-schema-language-go"
)

func main() error {
  // jsl.Schema can be parsed from JSON directly, but you can also construct
  // instances using the literal syntax:
  schema := jsl.Schema{
    RequiredProperties: map[string]jsl.Schema{
      "name":   jsl.Schema{Type: jsl.TypeString},
      "age":    jsl.Schema{Type: jsl.TypeNumber},
      "phones": jsl.Schema{
        Elements: &jsl.Schema{Type: jsl.TypeString}
      }
    }
  }

  // To keep this example simple, we'll construct this data by hand. But you
  // could also parse this data from JSON.
  //
  // This input data is perfect. It satisfies all the schema requirements.
  inputOk := map[string]interface{}{
    "name": "John Doe",
    "age":  43,
    "phones": [
      "+44 1234567",
      "+44 2345678",
    ],
  }

  // This input data has problems. "name" is missing, "age" has the wrong type,
  // and "phones[1]" has the wrong type.
  inputBad := map[string]interface{}{
    "age": "43",
    "phones": []interface{}{
      "+44 1234567",
      442345678,
    }
  }

  // To keep things simple, we'll ignore errors here. In this example, errors
  // are impossible. The docs explain in detail why an error might arise from
  // validation.
  validator := jsl.Validator{}
  resultOk, _ := validator.Validate(schema, inputOk)
  resultBad, _ := validator.Validate(schema, inputBad)

  fmt.Println(resultOk.IsValid()) // true
  fmt.Println(len(resultBad.Errors)) // 3

  // [] [properties name] -- indicates that the root is missing "name"
  fmt.Println(resultBad.Errors[0].InstancePath, resultBad.Errors[0].SchemaPath)

  // [age] [properties age type] -- indicates that "age" has the wrong type
  fmt.Println(resultBad.Errors[1].InstancePath, resultBad.Errors[1].SchemaPath)

  // [phones 1] [properties phones elements type] -- indicates that "phones[1]"
  // has the wrong type
  fmt.Println(resultBad.Errors[2].InstancePath, resultBad.Errors[2].SchemaPath)
}
```
