package jsl

import (
	"errors"
	"math"
	"strconv"
	"time"
)

type vm struct {
	MaxErrors               int
	MaxDepth                int
	StrictInstanceSemantics bool
	RootSchema              *Schema
	InstanceTokens          []string
	SchemaTokens            [][]string
	Errors                  []ValidationError
}

var errMaxErrors = errors.New("jsl internal: max errors reached")

func (vm *vm) validate(schema *Schema, instance *interface{}, parentTag *string) error {
	switch schema.Form() {
	case FormEmpty:
		// Nothing to be done. Empty never fails.
	case FormRef:
		if len(vm.SchemaTokens) == vm.MaxDepth {
			return MaxDepthExceededErr
		}

		refdSchema := vm.RootSchema.Definitions[*schema.Ref]
		vm.SchemaTokens = append(vm.SchemaTokens, []string{"definitions", *schema.Ref})

		if err := vm.validate(&refdSchema, instance, nil); err != nil {
			return err
		}

		vm.SchemaTokens = vm.SchemaTokens[:len(vm.SchemaTokens)-1]
	case FormType:
		switch schema.Type {
		case TypeBoolean:
			if _, ok := (*instance).(bool); !ok {
				vm.pushSchemaToken("type")
				if err := vm.pushErr(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		case TypeNumber, TypeFloat32, TypeFloat64:
			if _, ok := (*instance).(float64); !ok {
				vm.pushSchemaToken("type")
				if err := vm.pushErr(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		case TypeInt8:
			if err := vm.checkInt(instance, -128.0, 127.0); err != nil {
				return err
			}
		case TypeUint8:
			if err := vm.checkInt(instance, 0.0, 255.0); err != nil {
				return err
			}
		case TypeInt16:
			if err := vm.checkInt(instance, -32768.0, 32767.0); err != nil {
				return err
			}
		case TypeUint16:
			if err := vm.checkInt(instance, 0.0, 65535.0); err != nil {
				return err
			}
		case TypeInt32:
			if err := vm.checkInt(instance, -2147483648.0, 2147483647.0); err != nil {
				return err
			}
		case TypeUint32:
			if err := vm.checkInt(instance, 0.0, 4294967295.0); err != nil {
				return err
			}
		case TypeInt64:
			if err := vm.checkInt(instance, -9223372036854775808.0, 9223372036854775807.0); err != nil {
				return err
			}
		case TypeUint64:
			if err := vm.checkInt(instance, 0.0, 18446744073709551615.0); err != nil {
				return err
			}
		case TypeString:
			if _, ok := (*instance).(string); !ok {
				vm.pushSchemaToken("type")
				if err := vm.pushErr(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		case TypeTimestamp:
			if s, ok := (*instance).(string); ok {
				if _, err := time.Parse(time.RFC3339, s); err != nil {
					vm.pushSchemaToken("type")
					if err := vm.pushErr(); err != nil {
						return err
					}
					vm.popSchemaToken()
				}
			} else {
				vm.pushSchemaToken("type")
				if err := vm.pushErr(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		}
	case FormEnum:
		if s, ok := (*instance).(string); ok {
			found := false
			for _, val := range schema.Enum {
				if val == s {
					found = true
				}
			}

			if !found {
				vm.pushSchemaToken("enum")
				if err := vm.pushErr(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}
		} else {
			vm.pushSchemaToken("enum")
			if err := vm.pushErr(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case FormElements:
		if arr, ok := (*instance).([]interface{}); ok {
			vm.pushSchemaToken("elements")
			for i, elem := range arr {
				vm.pushInstanceToken(strconv.Itoa(i))
				if err := vm.validate(schema.Elements, &elem, nil); err != nil {
					return err
				}
				vm.popInstanceToken()
			}
			vm.popSchemaToken()
		} else {
			vm.pushSchemaToken("elements")
			if err := vm.pushErr(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case FormProperties:
		if obj, ok := (*instance).(map[string]interface{}); ok {
			vm.pushSchemaToken("properties")
			for property, subSchema := range schema.RequiredProperties {
				vm.pushSchemaToken(property)

				if val, ok := obj[property]; ok {
					vm.pushInstanceToken(property)
					if err := vm.validate(&subSchema, &val, nil); err != nil {
						return err
					}
					vm.popInstanceToken()
				} else {
					if err := vm.pushErr(); err != nil {
						return err
					}
				}

				vm.popSchemaToken()
			}
			vm.popSchemaToken()

			vm.pushSchemaToken("optionalProperties")
			for property, subSchema := range schema.OptionalProperties {
				vm.pushSchemaToken(property)

				if val, ok := obj[property]; ok {
					vm.pushInstanceToken(property)
					if err := vm.validate(&subSchema, &val, nil); err != nil {
						return err
					}
					vm.popInstanceToken()
				}

				vm.popSchemaToken()
			}
			vm.popSchemaToken()

			if vm.StrictInstanceSemantics {
				for k := range obj {
					// Do not apply strict instance semantics rule if the property is the
					// tag of a parent discriminator. This is the "discriminator tag
					// exemption" in the spec.
					if parentTag != nil && k == *parentTag {
						continue
					}

					requiredOk := false
					optionalOk := false

					if schema.RequiredProperties != nil {
						_, requiredOk = schema.RequiredProperties[k]
					}

					if schema.OptionalProperties != nil {
						_, optionalOk = schema.OptionalProperties[k]
					}

					if !requiredOk && !optionalOk {
						vm.pushInstanceToken(k)
						if err := vm.pushErr(); err != nil {
							return err
						}
						vm.popInstanceToken()
					}
				}
			}
		} else {
			// Sort of a weird corner-case in the spec: you have to check if the
			// instance is an object at all. If it isn't, you produce an error related
			// to `properties`. But if there wasn't a `properties` keyword, then you
			// have to produce `optionalProperties` instead.
			if schema.RequiredProperties != nil {
				vm.pushSchemaToken("properties")
			} else {
				vm.pushSchemaToken("optionalProperties")
			}

			if err := vm.pushErr(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case FormValues:
		if obj, ok := (*instance).(map[string]interface{}); ok {
			vm.pushSchemaToken("values")
			for k, v := range obj {
				vm.pushInstanceToken(k)
				if err := vm.validate(schema.Values, &v, nil); err != nil {
					return err
				}
				vm.popInstanceToken()
			}
			vm.popSchemaToken()
		} else {
			vm.pushSchemaToken("values")
			if err := vm.pushErr(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	case FormDiscriminator:
		if obj, ok := (*instance).(map[string]interface{}); ok {
			vm.pushSchemaToken("discriminator")

			if tagValue, ok := obj[schema.Discriminator.Tag]; ok {
				if tagValue, ok := tagValue.(string); ok {
					if subSchema, ok := schema.Discriminator.Mapping[tagValue]; ok {
						vm.pushSchemaToken("mapping")
						vm.pushSchemaToken(tagValue)
						if err := vm.validate(&subSchema, instance, &schema.Discriminator.Tag); err != nil {
							return err
						}
						vm.popSchemaToken()
						vm.popSchemaToken()
					} else {
						vm.pushSchemaToken("mapping")
						vm.pushInstanceToken(schema.Discriminator.Tag)
						if err := vm.pushErr(); err != nil {
							return err
						}
						vm.popInstanceToken()
						vm.popSchemaToken()
					}
				} else {
					vm.pushSchemaToken("tag")
					vm.pushInstanceToken(schema.Discriminator.Tag)
					if err := vm.pushErr(); err != nil {
						return err
					}
					vm.popInstanceToken()
					vm.popSchemaToken()
				}
			} else {
				vm.pushSchemaToken("tag")
				if err := vm.pushErr(); err != nil {
					return err
				}
				vm.popSchemaToken()
			}

			vm.popSchemaToken()
		} else {
			vm.pushSchemaToken("discriminator")
			if err := vm.pushErr(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	}

	return nil
}

func (vm *vm) checkInt(instance *interface{}, min, max float64) error {
	if n, ok := (*instance).(float64); ok {
		if i, f := math.Modf(n); f != 0.0 || i < min || i > max {
			vm.pushSchemaToken("type")
			if err := vm.pushErr(); err != nil {
				return err
			}
			vm.popSchemaToken()
		}
	} else {
		vm.pushSchemaToken("type")
		if err := vm.pushErr(); err != nil {
			return err
		}
		vm.popSchemaToken()
	}

	return nil
}

func (vm *vm) pushInstanceToken(token string) {
	vm.InstanceTokens = append(vm.InstanceTokens, token)
}

func (vm *vm) popInstanceToken() {
	vm.InstanceTokens = vm.InstanceTokens[:len(vm.InstanceTokens)-1]
}

func (vm *vm) pushSchemaToken(token string) {
	vm.SchemaTokens[len(vm.SchemaTokens)-1] = append(vm.SchemaTokens[len(vm.SchemaTokens)-1], token)
}

func (vm *vm) popSchemaToken() {
	schemaTokens := vm.SchemaTokens[len(vm.SchemaTokens)-1]
	vm.SchemaTokens[len(vm.SchemaTokens)-1] = schemaTokens[:len(schemaTokens)-1]
}

func (vm *vm) pushErr() error {
	instanceTokens := make([]string, len(vm.InstanceTokens))
	copy(instanceTokens, vm.InstanceTokens)

	schemaTokens := make([]string, len(vm.SchemaTokens[len(vm.SchemaTokens)-1]))
	copy(schemaTokens, vm.SchemaTokens[len(vm.SchemaTokens)-1])

	vm.Errors = append(vm.Errors, ValidationError{
		InstancePath: instanceTokens,
		SchemaPath:   schemaTokens,
	})

	if len(vm.Errors) == vm.MaxErrors {
		return errMaxErrors
	}

	return nil
}
