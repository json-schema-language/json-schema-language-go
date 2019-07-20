package jsl_test

import (
	"encoding/json"
	"testing"

	jsl "github.com/json-schema-language/json-schema-language-go"
	"github.com/stretchr/testify/assert"
)

func strptr(s string) *string {
	return &s
}

func TestVerifyAndForm(t *testing.T) {
	type testCase struct {
		in   string
		out  jsl.Schema
		err  error
		form jsl.Form
	}

	testCases := []testCase{
		{
			`{}`,
			jsl.Schema{},
			nil,
			jsl.FormEmpty,
		},
		{
			`{"ref":""}`,
			jsl.Schema{Ref: strptr("")},
			jsl.ErrNoSuchDefinition(""),
			jsl.FormRef,
		},
		{
			`{"definitions":{"":{}},"ref":""}`,
			jsl.Schema{
				Definitions: map[string]jsl.Schema{"": jsl.Schema{}},
				Ref:         strptr(""),
			},
			nil,
			jsl.FormRef,
		},
		{
			`{"definitions":{"":{}},"ref":"","type":"boolean"}`,
			jsl.Schema{
				Definitions: map[string]jsl.Schema{"": jsl.Schema{}},
				Ref:         strptr(""),
				Type:        jsl.TypeBoolean,
			},
			jsl.ErrInvalidForm,
			jsl.FormRef,
		},
		{
			`{"type":"boolean"}`,
			jsl.Schema{Type: jsl.TypeBoolean},
			nil,
			jsl.FormType,
		},
		{
			`{"type":"nonsense"}`,
			jsl.Schema{Type: "nonsense"},
			jsl.ErrInvalidType("nonsense"),
			jsl.FormType,
		},
		{
			`{"type":"boolean","enum":[]}`,
			jsl.Schema{Type: jsl.TypeBoolean, Enum: []string{}},
			jsl.ErrInvalidForm,
			jsl.FormType,
		},
		{
			`{"enum":["a","a"]}`,
			jsl.Schema{Enum: []string{"a", "a"}},
			jsl.ErrRepeatedEnumValue("a"),
			jsl.FormEnum,
		},
		{
			`{"enum":["a","b","c"]}`,
			jsl.Schema{Enum: []string{"a", "b", "c"}},
			nil,
			jsl.FormEnum,
		},
		{
			`{"enum":[],"properties":{}}`,
			jsl.Schema{Enum: []string{}, RequiredProperties: map[string]jsl.Schema{}},
			jsl.ErrInvalidForm,
			jsl.FormEnum,
		},
		{
			`{"enum":[],"elements":{}}`,
			jsl.Schema{Enum: []string{}, Elements: &jsl.Schema{}},
			jsl.ErrInvalidForm,
			jsl.FormEnum,
		},
		{
			`{"elements":{"ref":""}}`,
			jsl.Schema{Elements: &jsl.Schema{Ref: strptr("")}},
			jsl.ErrNoSuchDefinition(""),
			jsl.FormElements,
		},
		{
			`{"elements":{}}`,
			jsl.Schema{Elements: &jsl.Schema{}},
			nil,
			jsl.FormElements,
		},
		{
			`{"elements":{},"properties":{}}`,
			jsl.Schema{Elements: &jsl.Schema{}, RequiredProperties: map[string]jsl.Schema{}},
			jsl.ErrInvalidForm,
			jsl.FormElements,
		},
		{
			`{"elements":{},"optionalProperties":{}}`,
			jsl.Schema{Elements: &jsl.Schema{}, OptionalProperties: map[string]jsl.Schema{}},
			jsl.ErrInvalidForm,
			jsl.FormElements,
		},
		{
			`{"properties":{"a":{}},"optionalProperties":{"a":{}}}`,
			jsl.Schema{
				RequiredProperties: map[string]jsl.Schema{"a": jsl.Schema{}},
				OptionalProperties: map[string]jsl.Schema{"a": jsl.Schema{}},
			},
			jsl.ErrRepeatedProperty("a"),
			jsl.FormProperties,
		},
		{
			`{"properties":{"a":{"ref":""}}}`,
			jsl.Schema{
				RequiredProperties: map[string]jsl.Schema{"a": jsl.Schema{Ref: strptr("")}},
			},
			jsl.ErrNoSuchDefinition(""),
			jsl.FormProperties,
		},
		{
			`{"optionalProperties":{"a":{"ref":""}}}`,
			jsl.Schema{
				OptionalProperties: map[string]jsl.Schema{"a": jsl.Schema{Ref: strptr("")}},
			},
			jsl.ErrNoSuchDefinition(""),
			jsl.FormProperties,
		},
		{
			`{"properties":{"a":{}},"optionalProperties":{"b":{}}}`,
			jsl.Schema{
				RequiredProperties: map[string]jsl.Schema{"a": jsl.Schema{}},
				OptionalProperties: map[string]jsl.Schema{"b": jsl.Schema{}},
			},
			nil,
			jsl.FormProperties,
		},
		{
			`{"properties":{},"values":{}}`,
			jsl.Schema{
				RequiredProperties: map[string]jsl.Schema{},
				Values:             &jsl.Schema{},
			},
			jsl.ErrInvalidForm,
			jsl.FormProperties,
		},
		{
			`{"values":{"ref":""}}`,
			jsl.Schema{Values: &jsl.Schema{Ref: strptr("")}},
			jsl.ErrNoSuchDefinition(""),
			jsl.FormValues,
		},
		{
			`{"values":{}}`,
			jsl.Schema{Values: &jsl.Schema{}},
			nil,
			jsl.FormValues,
		},
		{
			`{"values":{},"discriminator":{"mapping":{}}}`,
			jsl.Schema{
				Values:        &jsl.Schema{},
				Discriminator: jsl.Discriminator{Mapping: map[string]jsl.Schema{}},
			},
			jsl.ErrInvalidForm,
			jsl.FormValues,
		},
		{
			`{"discriminator":{"mapping":{"":{}}}}`,
			jsl.Schema{
				Discriminator: jsl.Discriminator{Mapping: map[string]jsl.Schema{
					"": jsl.Schema{},
				}},
			},
			jsl.ErrNonPropertiesMapping,
			jsl.FormDiscriminator,
		},
		{
			`{"discriminator":{"tag":"a","mapping":{"":{"properties":{"a":{}}}}}}`,
			jsl.Schema{
				Discriminator: jsl.Discriminator{
					Tag: "a",
					Mapping: map[string]jsl.Schema{
						"": jsl.Schema{
							RequiredProperties: map[string]jsl.Schema{
								"a": jsl.Schema{},
							},
						},
					},
				},
			},
			jsl.ErrRepeatedTagInProperties("a"),
			jsl.FormDiscriminator,
		},
		{
			`{"discriminator":{"tag":"a","mapping":{"":{"optionalProperties":{"a":{}}}}}}`,
			jsl.Schema{
				Discriminator: jsl.Discriminator{
					Tag: "a",
					Mapping: map[string]jsl.Schema{
						"": jsl.Schema{
							OptionalProperties: map[string]jsl.Schema{
								"a": jsl.Schema{},
							},
						},
					},
				},
			},
			jsl.ErrRepeatedTagInProperties("a"),
			jsl.FormDiscriminator,
		},
		{
			`{"discriminator":{"tag":"a","mapping":{"":{"properties":{"b":{}}}}}}`,
			jsl.Schema{
				Discriminator: jsl.Discriminator{
					Tag: "a",
					Mapping: map[string]jsl.Schema{
						"": jsl.Schema{
							RequiredProperties: map[string]jsl.Schema{
								"b": jsl.Schema{},
							},
						},
					},
				},
			},
			nil,
			jsl.FormDiscriminator,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.in, func(t *testing.T) {
			var out jsl.Schema
			assert.NoError(t, json.Unmarshal([]byte(tt.in), &out))

			assert.Equal(t, tt.out, out)
			assert.Equal(t, tt.err, out.Verify())
			assert.Equal(t, tt.form, out.Form())
		})
	}
}
