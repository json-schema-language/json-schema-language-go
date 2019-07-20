package jsl

import (
	"errors"
	"fmt"
)

// ErrInvalidForm indicates that the schema does not fall into one of the eight
// forms.
var ErrInvalidForm = errors.New("jsl: ambiguous or invalid schema form")

// ErrNonPropertiesMapping indicates that the schema had a Discriminator.Mapping
// containing schemas that weren't of the properties form.
//
// Per the spec, all discriminator mapping values must be of the properties
// form, or the schema is not correct.
var ErrNonPropertiesMapping = errors.New("jsl: value of discriminator mapping is not of properties form")

// ErrMaxDepthExceeded indicates that the maximum evaluation depth was exceeded
// while validating an instance. This typically indicates that an infinite
// recurisve loop was encountered while evaluating the schema.
var ErrMaxDepthExceeded = errors.New("jsl: maximum evaluation depth exceeded")

// ErrNoSuchDefinition indicates that a "ref" referred to a definition that does
// not exist.
type ErrNoSuchDefinition string

func (e ErrNoSuchDefinition) Error() string {
	return fmt.Sprintf("jsl: no such definition: %s", e)
}

// ErrInvalidType indicates that a "type" had an incorrect value.
//
// See Type for all the correct values "type" may take on.
type ErrInvalidType string

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("jsl: no such type: %s", e)
}

// ErrRepeatedEnumValue indicates than an "enum" repeated a value. Enums must
// not contain duplicates.
type ErrRepeatedEnumValue string

func (e ErrRepeatedEnumValue) Error() string {
	return fmt.Sprintf("jsl: repeated enum value: %s", e)
}

// ErrRepeatedProperty indicates that a schema had a "properties" and
// "optionalProperties" that specified the same property.
type ErrRepeatedProperty string

func (e ErrRepeatedProperty) Error() string {
	return fmt.Sprintf("jsl: repeated property in properties and optionalProperties: %s", e)
}

// ErrRepeatedTagInProperties indicates that one of the elements of
// Discriminator.Mapping repeated the Discriminator.Tag in one of its
// "properties" or "optionalProperties".
type ErrRepeatedTagInProperties string

func (e ErrRepeatedTagInProperties) Error() string {
	return fmt.Sprintf("jsl: discriminator tag repeated in properties or optionalProperties: %s", e)
}
