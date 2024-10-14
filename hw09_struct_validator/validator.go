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
	ErrFailNotStruct        = errors.New("input interface is not a struct")
	ErrFailNotSupportedType = errors.New("validator unsupported type")
)

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrFailNotStruct
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
		var err error
		validator.SetFields(tag, field.Name, fieldValue)
		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.String:
			validator.SetStringValidator()
			err = validator.Validate()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			validator.SetIntegerValidator()
			err = validator.Validate()
		case reflect.Slice:
			switch fieldValue.Type().String() {
			case "[]string":
				validator.SetStringValidator()
				err = validateSlices(validator, fieldValue)
			case "[]int", "[]int8", "[]int16", "[]in32", "[]int64":
				validator.SetIntegerValidator()
				err = validateSlices(validator, fieldValue)
			}
		case reflect.Struct:
			if tag == "nested" {
				err = Validate(fieldValue.Interface())
				validator.wrapErrors(err)
			}
		default:
			return ErrFailNotSupportedType
		}
		if err != nil {
			return err
		}
	}

	if len(validator.errs) > 0 {
		return validator.errs
	}
	return nil
}

func validateSlices(validator *Validator, fieldValue reflect.Value) error {
	var err error
	for i := 0; i < fieldValue.Len(); i++ {
		validator.SetReflectValue(fieldValue.Index(i))
		err = validator.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
