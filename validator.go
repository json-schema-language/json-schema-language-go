package jsl

type Validator struct {
	MaxErrors               int
	MaxDepth                int
	StrictInstanceSemantics bool
}

type ValidationResult struct {
	Errors []ValidationError
}

type ValidationError struct {
	InstancePath []string
	SchemaPath   []string
}

func (v *Validator) Validate(schema *Schema, instance interface{}) (ValidationResult, error) {
	vm := vm{
		MaxErrors:               v.MaxErrors,
		MaxDepth:                v.MaxDepth,
		StrictInstanceSemantics: v.StrictInstanceSemantics,
		RootSchema:              schema,
		InstanceTokens:          []string{},
		SchemaTokens:            [][]string{[]string{}},
	}

	if err := vm.validate(schema, &instance, nil); err != nil && err != errMaxErrors {
		return ValidationResult{}, err
	}

	return ValidationResult{Errors: vm.Errors}, nil
}
