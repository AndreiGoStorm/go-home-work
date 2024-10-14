package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrFailWrongMinNumber = errors.New("value for min is not a number")
	ErrFailWrongMaxNumber = errors.New("value for max is not a number")
	ErrFailSetNumber      = errors.New("value for in is not a number")
)

var (
	ErrIntLessThanMin      = errors.New("less than minimum")
	ErrIntMoreThanMax      = errors.New("more than maximum")
	ErrIntNotIncludedInSet = errors.New("not included in validation set")
)

type IntegerValidator struct {
	v     *Validator
	value int64
}

func (iv *IntegerValidator) prepareParams(v *Validator) {
	iv.v = v
	iv.value = v.value.Int()
}

func (iv *IntegerValidator) validate(v *Validator) error {
	iv.prepareParams(v)
	for _, tag := range v.tags {
		args, err := v.parseTag(tag)
		if err != nil {
			return err
		}

		switch args[0] {
		case "min":
			err = iv.Min(args[1])
		case "max":
			err = iv.Max(args[1])
		case "in":
			err = iv.In(args[1])
		default:
			return fmt.Errorf("%w: %s, for tag: %s", ErrFailWrongRule, v.field, tag)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (iv *IntegerValidator) Min(arg string) error {
	condition, err := strconv.Atoi(arg)
	if err != nil {
		return ErrFailWrongMinNumber
	}
	if int64(condition) > iv.value {
		iv.v.addValidationError(fmt.Errorf("%w: condition %d, value %d", ErrIntLessThanMin, condition, iv.value))
	}
	return nil
}

func (iv *IntegerValidator) Max(arg string) error {
	condition, err := strconv.Atoi(arg)
	if err != nil {
		return ErrFailWrongMaxNumber
	}
	if int64(condition) < iv.value {
		iv.v.addValidationError(fmt.Errorf("%w: condition %d, value %d", ErrIntMoreThanMax, condition, iv.value))
	}
	return nil
}

func (iv *IntegerValidator) In(arg string) error {
	set := strings.Split(arg, ",")
	inSet := false
	for _, s := range set {
		intVal, err := strconv.Atoi(s)
		if err != nil {
			return ErrFailSetNumber
		}
		if int64(intVal) == iv.value {
			inSet = true
			break
		}
	}
	if !inSet {
		iv.v.addValidationError(fmt.Errorf("%w: value %d, set %v", ErrIntNotIncludedInSet, iv.value, set))
	}
	return nil
}
