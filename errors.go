package jsl

import (
	"errors"
	"fmt"
)

var InvalidFormErr = errors.New("jsl: ambiguous or invalid schema form")
var NonPropertiesMappingErr = errors.New("jsl: value of discriminator mapping is not of properties form")
var MaxDepthExceededErr = errors.New("jsl: maximum evaluation depth exceeded")

type NoSuchDefinitionErr string

func (e NoSuchDefinitionErr) Error() string {
	return fmt.Sprintf("jsl: no such definition: %s", e)
}

type InvalidTypeErr string

func (e InvalidTypeErr) Error() string {
	return fmt.Sprintf("jsl: no such type: %s", e)
}

type RepeatedEnumValueErr string

func (e RepeatedEnumValueErr) Error() string {
	return fmt.Sprintf("jsl: repeated enum value: %s", e)
}

type RepeatedPropertyErr string

func (e RepeatedPropertyErr) Error() string {
	return fmt.Sprintf("jsl: repeated property in properties and optionalProperties: %s", e)
}

type RepeatedTagInPropertiesErr string

func (e RepeatedTagInPropertiesErr) Error() string {
	return fmt.Sprintf("jsl: discriminator tag repeated in properties or optionalProperties: %s", e)
}
