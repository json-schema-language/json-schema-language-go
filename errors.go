package jsl

import (
	"errors"
	"fmt"
)

var ErrInvalidForm = errors.New("jsl: ambiguous or invalid schema form")
var ErrNonPropertiesMapping = errors.New("jsl: value of discriminator mapping is not of properties form")
var ErrMaxDepthExceeded = errors.New("jsl: maximum evaluation depth exceeded")

type ErrNoSuchDefinition string

func (e ErrNoSuchDefinition) Error() string {
	return fmt.Sprintf("jsl: no such definition: %s", e)
}

type ErrInvalidType string

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("jsl: no such type: %s", e)
}

type ErrRepeatedEnumValue string

func (e ErrRepeatedEnumValue) Error() string {
	return fmt.Sprintf("jsl: repeated enum value: %s", e)
}

type ErrRepeatedProperty string

func (e ErrRepeatedProperty) Error() string {
	return fmt.Sprintf("jsl: repeated property in properties and optionalProperties: %s", e)
}

type ErrRepeatedTagInProperties string

func (e ErrRepeatedTagInProperties) Error() string {
	return fmt.Sprintf("jsl: discriminator tag repeated in properties or optionalProperties: %s", e)
}
