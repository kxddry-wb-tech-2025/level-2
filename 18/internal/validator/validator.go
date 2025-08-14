package validator

import "github.com/go-playground/validator/v10"

// Validator is a struct for echo middleware to validate incoming events
type Validator struct {
	val *validator.Validate
}

// New creates a Validator
func New(val *validator.Validate) *Validator {
	return &Validator{val: val}
}

// Validate validates an event
func (v *Validator) Validate(i interface{}) error {
	return v.val.Struct(i)
}
