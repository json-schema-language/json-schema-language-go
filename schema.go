package jsl

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

type Form int

const (
	FormEmpty Form = iota
	FormRef
	FormType
	FormEnum
	FormElements
	FormProperties
	FormValues
	FormDiscriminator
)

type Type string

const (
	TypeBoolean   Type = "boolean"
	TypeNumber         = "number"
	TypeFloat32        = "float32"
	TypeFloat64        = "float64"
	TypeInt8           = "int8"
	TypeUint8          = "uint8"
	TypeInt16          = "int16"
	TypeUint16         = "uint16"
	TypeInt32          = "int32"
	TypeUint32         = "uint32"
	TypeInt64          = "int64"
	TypeUint64         = "uint64"
	TypeString         = "string"
	TypeTimestamp      = "timestamp"
)

type Discriminator struct {
	Tag     string            `json:"tag"`
	Mapping map[string]Schema `json:"mapping"`
}

func (s *Schema) Verify() error {
	return s.verify(s)
}

func (s *Schema) verify(root *Schema) error {
	isEmpty := true

	if s.Ref != nil {
		if root.Definitions == nil {
			return NoSuchDefinitionErr(*s.Ref)
		}

		if _, ok := root.Definitions[*s.Ref]; !ok {
			return NoSuchDefinitionErr(*s.Ref)
		}

		isEmpty = false
	}

	if s.Type != "" {
		if !isEmpty {
			return InvalidFormErr
		}

		switch s.Type {
		case "boolean", "number", "float32", "float64", "int8", "uint8", "int16",
			"uint16", "int32", "uint32", "int64", "uint64", "string", "timestamp":
		default:
			return InvalidTypeErr(s.Type)
		}

		isEmpty = false
	}

	if s.Enum != nil {
		if !isEmpty {
			return InvalidFormErr
		}

		vals := map[string]struct{}{}
		for _, val := range s.Enum {
			if _, ok := vals[val]; ok {
				return RepeatedEnumValueErr(val)
			}

			vals[val] = struct{}{}
		}

		isEmpty = false
	}

	if s.Elements != nil {
		if !isEmpty {
			return InvalidFormErr
		}

		if err := s.Elements.verify(root); err != nil {
			return err
		}

		isEmpty = false
	}

	if s.RequiredProperties != nil || s.OptionalProperties != nil {
		if !isEmpty {
			return InvalidFormErr
		}

		if s.RequiredProperties != nil && s.OptionalProperties != nil {
			for k := range s.RequiredProperties {
				if _, ok := s.OptionalProperties[k]; ok {
					return RepeatedPropertyErr(k)
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
			return InvalidFormErr
		}

		if err := s.Values.verify(root); err != nil {
			return err
		}

		isEmpty = false
	}

	if s.Discriminator.Mapping != nil {
		if !isEmpty {
			return InvalidFormErr
		}

		for _, m := range s.Discriminator.Mapping {
			if err := m.verify(root); err != nil {
				return err
			}

			if m.Form() != FormProperties {
				return NonPropertiesMappingErr
			}

			for k := range m.RequiredProperties {
				if k == s.Discriminator.Tag {
					return RepeatedTagInPropertiesErr(s.Discriminator.Tag)
				}
			}

			for k := range m.OptionalProperties {
				if k == s.Discriminator.Tag {
					return RepeatedTagInPropertiesErr(s.Discriminator.Tag)
				}
			}
		}

		isEmpty = false
	}

	return nil
}

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
