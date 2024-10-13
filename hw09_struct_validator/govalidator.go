package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrValidationParse = errors.New("parsing error in field")
	ErrUndefinedRule   = errors.New("undefined rule in field")
)

type iValidator interface {
	validate(v *Validator)
}

type Validator struct {
	iValidator iValidator
	tags       []string
	value      reflect.Value
	field      string
	errs       ValidationErrors
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) SetFields(tag, field string, value reflect.Value) {
	v.tags = strings.Split(tag, "|")
	v.field = field
	v.value = value
}

func (v *Validator) SetReflectValue(value reflect.Value) {
	v.value = value
}

func (v *Validator) SetIValidator(i iValidator) {
	v.iValidator = i
}

func (v *Validator) SetStringValidator() {
	v.SetIValidator(&StringValidator{})
}

func (v *Validator) SetIntegerValidator() {
	v.SetIValidator(&IntegerValidator{})
}

func (v *Validator) Validate() {
	v.iValidator.validate(v)
}

func (v *Validator) parseTag(tag string) []string {
	args := strings.Split(tag, ":")
	if len(args) < 2 {
		v.addValidationError(fmt.Errorf("%w: %s, for tag: %s", ErrValidationParse, v.field, tag))
		return nil
	}

	return args
}

func (v *Validator) addValidationError(err error) {
	v.errs = append(v.errs, ValidationError{v.field, err})
}

func (v *Validator) wrapErrors(errs error) {
	var validationErrs ValidationErrors
	if errors.As(errs, &validationErrs) {
		v.errs = append(v.errs, validationErrs...)
	}
}
