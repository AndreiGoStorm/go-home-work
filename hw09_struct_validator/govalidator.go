package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrFailWrongArgs = errors.New("wrong arguments in field")
	ErrFailWrongRule = errors.New("wrong rule in field")
)

type iValidator interface {
	validate(v *Validator) error
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

func (v *Validator) Validate() error {
	return v.iValidator.validate(v)
}

func (v *Validator) parseTag(tag string) ([]string, error) {
	args := strings.Split(tag, ":")
	if len(args) < 2 {
		return nil, fmt.Errorf("%w: %s, for tag: %s", ErrFailWrongArgs, v.field, tag)
	}

	return args, nil
}

func (v *Validator) addValidationError(err error) {
	v.errs = append(v.errs, ValidationError{v.field, err})
}

func (v *Validator) wrapErrors(err error) {
	var validationErrs ValidationErrors
	if errors.As(err, &validationErrs) {
		v.errs = append(v.errs, validationErrs...)
	}
}
