package jsl

// Validator validates instances ("inputs") against schemas.
//
// When evaluating untrusted schemas, always set MaxDepth to a nonzero value.
type Validator struct {
	// The maximum number of errors to produce from Validate. Zero indicates that
	// all errors should be produced.
	MaxErrors int

	// The maximum stack depth when evaluating schemas. Schemas may potentially
	// have circular recursive definitions using the "ref" keyword, and so it's
	// possible to encounter stack overflows when calling Validate. By setting
	// MaxDepth to a positive value, Validator can abort validation with an error
	// when a certain number of "ref" calls have been recursively followed.
	//
	// When evaluating schemas, always set this to a nonzero value, or you may
	// overflow the stack.
	//
	// Zero indicates that no maximum depth should be imposed.
	MaxDepth int

	// Whether to enforce strict instance semantics. See the spec for a formal
	// definition, but essentially, strict instance semantics bans "unknown" or
	// "unspecified" properties from appearing in instances.
	StrictInstanceSemantics bool
}

// ValidationResult is the set of validation errors arising from running
// Validate.
type ValidationResult struct {
	Errors []ValidationError
}

// IsValid returns whether a ValidationResult indicates that the instance is
// valid against the schema.
//
// This just checks whether Errors is empty.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// ValidationError is a single validation error. InstancePath and SchemaPath are
// the tokens of a JSON Pointer. They point to the part of the instance which
// had a problem, and the part of the schema which reported the problem.
//
// See the spec for a formal definition of how InstancePath and SchemaPath are
// computed.
//
// Note that this is not a error in the Golang sense of the term.
type ValidationError struct {
	InstancePath []string
	SchemaPath   []string
}

// Validate checks whether an instance ("input") is valid against a Schema, and
// reports the validation errors that arose while doing this check.
//
// If the instance is invalid, this does not cause an error to be returned.
// ValidationResult will contain all validation errors, and is not an error in
// the Golang sense of the term.
//
// ErrMaxDepthExceeded is returned if the maximum depth is exceeded. See
// MaxDepth on Validator for more details.
func (v *Validator) Validate(schema Schema, instance interface{}) (ValidationResult, error) {
	vm := vm{
		MaxErrors:               v.MaxErrors,
		MaxDepth:                v.MaxDepth,
		StrictInstanceSemantics: v.StrictInstanceSemantics,
		RootSchema:              schema,
		InstanceTokens:          []string{},
		SchemaTokens:            [][]string{[]string{}},
	}

	if err := vm.validate(schema, instance, nil); err != nil && err != errMaxErrors {
		return ValidationResult{}, err
	}

	return ValidationResult{Errors: vm.Errors}, nil
}
