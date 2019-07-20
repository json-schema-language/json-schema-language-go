package jsl_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/dolmen-go/jsonptr"
	jsl "github.com/json-schema-language/json-schema-language-go"
	"github.com/stretchr/testify/assert"
)

type validationError struct {
	InstancePath string `json:"instancePath"`
	SchemaPath   string `json:"schemaPath"`
}

type instanceCase struct {
	Instance interface{}       `json:"instance"`
	Errors   []validationError `json:"errors"`
}

type testCase struct {
	Name           string         `json:"name"`
	Schema         jsl.Schema     `json:"schema"`
	StrictInstance bool           `json:"strictInstance"`
	Instances      []instanceCase `json:"instances"`
}

func sortErrors(errs []validationError) {
	sort.Slice(errs, func(i, j int) bool {
		if errs[i].SchemaPath == errs[j].SchemaPath {
			return errs[i].InstancePath < errs[j].InstancePath
		}

		return errs[i].SchemaPath < errs[j].SchemaPath
	})
}

func TestSpec(t *testing.T) {
	assert.NoError(t,
		filepath.Walk("spec/tests", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			var testCases []testCase
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&testCases); err != nil {
				return err
			}

			for _, tt := range testCases {
				t.Run(fmt.Sprintf("%s/%s", path, tt.Name), func(t *testing.T) {
					validator := jsl.Validator{StrictInstanceSemantics: tt.StrictInstance}

					for i, instance := range tt.Instances {
						t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
							result, err := validator.Validate(&tt.Schema, &instance.Instance)
							assert.NoError(t, err)

							// Stringify result's errors into JSON Pointers for comparison
							// with spec test cases.
							errors := make([]validationError, len(result.Errors))
							for i, err := range result.Errors {
								errors[i] = validationError{
									InstancePath: jsonptr.Pointer(err.InstancePath).String(),
									SchemaPath:   jsonptr.Pointer(err.SchemaPath).String(),
								}
							}

							sortErrors(instance.Errors)
							sortErrors(errors)
							assert.Equal(t, instance.Errors, errors)
						})
					}
				})
			}

			return nil
		}))
}
