package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrIntMinNumber        = errors.New("value for min is not a number")
	ErrIntLessThanMin      = errors.New("less than minimum")
	ErrIntMaxNumber        = errors.New("value for max is not a number")
	ErrIntMoreThanMax      = errors.New("more than maximum")
	ErrIntSetNumber        = errors.New("value for in is not a number")
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

func (iv *IntegerValidator) validate(v *Validator) {
	iv.prepareParams(v)
	for _, tag := range v.tags {
		args := v.parseTag(tag)
		if args == nil {
			return
		}

		switch args[0] {
		case "min":
			iv.Min(args[1])
		case "max":
			iv.Max(args[1])
		case "in":
			iv.In(args[1])
		default:
			iv.v.addValidationError(fmt.Errorf("%w: %s, for tag: %s", ErrUndefinedRule, v.field, tag))
		}
	}
}

func (iv *IntegerValidator) Min(arg string) {
	condition, err := strconv.Atoi(arg)
	if err != nil {
		iv.v.addValidationError(ErrIntMinNumber)
		return
	}
	if int64(condition) > iv.value {
		iv.v.addValidationError(fmt.Errorf("%w: condition %d, value %d", ErrIntLessThanMin, condition, iv.value))
	}
}

func (iv *IntegerValidator) Max(arg string) {
	condition, err := strconv.Atoi(arg)
	if err != nil {
		iv.v.addValidationError(ErrIntMaxNumber)
		return
	}
	if int64(condition) < iv.value {
		iv.v.addValidationError(fmt.Errorf("%w: condition %d, value %d", ErrIntMoreThanMax, condition, iv.value))
	}
}

func (iv *IntegerValidator) In(arg string) {
	set := strings.Split(arg, ",")
	inSet := false
	for _, s := range set {
		intVal, err := strconv.Atoi(s)
		if err != nil {
			iv.v.addValidationError(ErrIntSetNumber)
			return
		}
		if int64(intVal) == iv.value {
			inSet = true
			break
		}
	}
	if !inSet {
		iv.v.addValidationError(fmt.Errorf("%w: value %d, set %v", ErrIntNotIncludedInSet, iv.value, set))
	}
}
