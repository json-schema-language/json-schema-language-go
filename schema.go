package jsl

// Schema represents a JSON Schema Language schema.
//
// This type is designed for conversion to/from JSON. However, not all instances
// of this type are "correct" schemas as defined by the JSL spec. To verify the
// correctness of a schema, use Verify.
type Schema struct {
	Definitions        map[string]Schema `json:"definitions"`
	Ref                *string           `json:"ref"`
	Type               Type              `json:"type"`
	Enum               []string          `json:"enum"`
	Elements           *Schema           `json:"elements"`
	RequiredProperties map[string]Schema `json:"properties"`
	OptionalProperties map[string]Schema `json:"optionalProperties"`
	Values             *Schema           `json:"values"`
	Discriminator      Discriminator     `json:"discriminator"`
}

// Type represents the correct values for Type in Schema.
type Type string

const (
	// TypeBoolean represents true and false.
	TypeBoolean Type = "boolean"

	// TypeNumber represents all numbers.
	TypeNumber = "number"

	// TypeFloat32 represents a float32.
	TypeFloat32 = "float32"

	// TypeFloat64 represents a float64.
	TypeFloat64 = "float64"

	// TypeInt8 represents a int8.
	TypeInt8 = "int8"

	// TypeUint8 represents a uint8.
	TypeUint8 = "uint8"

	// TypeInt16 represents a int16.
	TypeInt16 = "int16"

	// TypeUint16 represents a uint16.
	TypeUint16 = "uint16"

	// TypeInt32 represents a int32.
	TypeInt32 = "int32"

	// TypeUint32 represents a uint32.
	TypeUint32 = "uint32"

	// TypeInt64 represents a int64.
	TypeInt64 = "int64"

	// TypeUint64 represents a uint64.
	TypeUint64 = "uint64"

	// TypeString represents a string.
	TypeString = "string"

	// TypeTimestamp represents a string encoding a RFC3339 timestamp.
	TypeTimestamp = "timestamp"
)

// Discriminator stores data associated with a schema of the discriminator form.
type Discriminator struct {
	Tag     string            `json:"tag"`
	Mapping map[string]Schema `json:"mapping"`
}

// Form represents the eight kinds of JSL schemas defined by the spec.
//
// All correct schemas conform to exactly one of the eight forms.
type Form int

const (
	// FormEmpty represents the "empty" form.
	FormEmpty Form = iota

	// FormRef represents the "ref" form.
	FormRef

	// FormType represents the "type" form.
	FormType

	// FormEnum represents the "enum" form.
	FormEnum

	// FormElements represents the "elements" form.
	FormElements

	// FormProperties represents the "properties" form.
	FormProperties

	// FormValues represents the "values" form.
	FormValues

	// FormDiscriminator represents the "discriminator" form.
	FormDiscriminator
)

// Form determines which form a schema takes on, assuming it is correct.
//
// If the Schema is not correct, then this function's return value is not
// meaningful.
func (s *Schema) Form() Form {
	if s.Ref != nil {
		return FormRef
	} else if s.Type != "" {
		return FormType
	} else if s.Enum != nil {
		return FormEnum
	} else if s.Elements != nil {
		return FormElements
	} else if s.RequiredProperties != nil || s.OptionalProperties != nil {
		return FormProperties
	} else if s.Values != nil {
		return FormValues
	} else if s.Discriminator.Mapping != nil {
		return FormDiscriminator
	} else {
		return FormEmpty
	}
}

// Verify returns nil if a schema is correct, or an error if it is not. The
// error contains details on the first encountered problem with the correctness
// of the schema.
func (s *Schema) Verify() error {
	return s.verify(s)
}

func (s *Schema) verify(root *Schema) error {
	isEmpty := true

	if s.Ref != nil {
		if root.Definitions == nil {
			return ErrNoSuchDefinition(*s.Ref)
		}

		if _, ok := root.Definitions[*s.Ref]; !ok {
			return ErrNoSuchDefinition(*s.Ref)
		}

		isEmpty = false
	}

	if s.Type != "" {
		if !isEmpty {
			return ErrInvalidForm
		}

		switch s.Type {
		case "boolean", "number", "float32", "float64", "int8", "uint8", "int16",
			"uint16", "int32", "uint32", "int64", "uint64", "string", "timestamp":
		default:
			return ErrInvalidType(s.Type)
		}

		isEmpty = false
	}

	if s.Enum != nil {
		if !isEmpty {
			return ErrInvalidForm
		}

		vals := map[string]struct{}{}
		for _, val := range s.Enum {
			if _, ok := vals[val]; ok {
				return ErrRepeatedEnumValue(val)
			}

			vals[val] = struct{}{}
		}

		isEmpty = false
	}

	if s.Elements != nil {
		if !isEmpty {
			return ErrInvalidForm
		}

		if err := s.Elements.verify(root); err != nil {
			return err
		}

		isEmpty = false
	}

	if s.RequiredProperties != nil || s.OptionalProperties != nil {
		if !isEmpty {
			return ErrInvalidForm
		}

		if s.RequiredProperties != nil && s.OptionalProperties != nil {
			for k := range s.RequiredProperties {
				if _, ok := s.OptionalProperties[k]; ok {
					return ErrRepeatedProperty(k)
				}
			}
		}

		for _, s := range s.RequiredProperties {
			if err := s.verify(root); err != nil {
				return err
			}
		}

		for _, s := range s.OptionalProperties {
			if err := s.verify(root); err != nil {
				return err
			}
		}

		isEmpty = false
	}

	if s.Values != nil {
		if !isEmpty {
			return ErrInvalidForm
		}

		if err := s.Values.verify(root); err != nil {
			return err
		}

		isEmpty = false
	}

	if s.Discriminator.Mapping != nil {
		if !isEmpty {
			return ErrInvalidForm
		}

		for _, m := range s.Discriminator.Mapping {
			if err := m.verify(root); err != nil {
				return err
			}

			if m.Form() != FormProperties {
				return ErrNonPropertiesMapping
			}

			for k := range m.RequiredProperties {
				if k == s.Discriminator.Tag {
					return ErrRepeatedTagInProperties(s.Discriminator.Tag)
				}
			}

			for k := range m.OptionalProperties {
				if k == s.Discriminator.Tag {
					return ErrRepeatedTagInProperties(s.Discriminator.Tag)
				}
			}
		}

		isEmpty = false
	}

	return nil
}
