package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for index, err := range v {
		builder.WriteString("Error in field `")
		builder.WriteString(err.Field)
		builder.WriteString("` ")
		builder.WriteString(err.Err.Error())
		builder.WriteString(".")
		if index != len(v)-1 {
			builder.WriteString(" ")
		}
	}
	return builder.String()
}

var (
	ErrNotStruct        = errors.New("input interface is not a struct")
	ErrNotSupportedType = errors.New("validator unsupported type")
)

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	validator := NewValidator()
	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		tag := field.Tag.Get("validate")
		if len(tag) == 0 {
			continue
		}

		fieldValue := value.Field(i)
		if !fieldValue.CanInterface() {
			continue
		}

		validator.SetFields(tag, field.Name, fieldValue)
		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.String:
			validator.SetStringValidator()
			validator.Validate()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			validator.SetIntegerValidator()
			validator.Validate()
		case reflect.Slice:
			switch fieldValue.Type().String() {
			case "[]string":
				validator.SetStringValidator()
				validateSlices(validator, fieldValue)
			case "[]int", "[]int8", "[]int16", "[]in32", "[]int64":
				validator.SetIntegerValidator()
				validateSlices(validator, fieldValue)
			}
		case reflect.Struct:
			if tag == "nested" {
				errs := Validate(fieldValue.Interface())
				validator.wrapErrors(errs)
			}
		default:
			return ErrNotSupportedType
		}
	}

	if len(validator.errs) > 0 {
		return validator.errs
	}
	return nil
}

func validateSlices(validator *Validator, fieldValue reflect.Value) {
	for i := 0; i < fieldValue.Len(); i++ {
		validator.SetReflectValue(fieldValue.Index(i))
		validator.Validate()
	}
}
