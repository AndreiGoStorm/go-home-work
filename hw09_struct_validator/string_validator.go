package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrFailWrongLengthNumber = errors.New("value for length is not a number")
	ErrFailWrongRegexp       = errors.New("regexp is wrong")
	ErrFailWrongSet          = errors.New("wrong parameter for in")
)

var (
	ErrStringLength           = errors.New("length of string is wrong")
	ErrStringNotMatchedRegexp = errors.New("string is not matched regexp")
	ErrStringNotIncludedSet   = errors.New("string is not included in set")
)

type StringValidator struct {
	v     *Validator
	value string
}

func (sv *StringValidator) prepareParams(v *Validator) {
	sv.v = v
	sv.value = v.value.String()
}

func (sv *StringValidator) validate(v *Validator) error {
	sv.prepareParams(v)
	for _, tag := range v.tags {
		args, err := v.parseTag(tag)
		if err != nil {
			return err
		}
		switch args[0] {
		case "len":
			err = sv.Len(args[1])
		case "regexp":
			err = sv.Regexp(args[1])
		case "in":
			err = sv.In(args[1])
		default:
			return fmt.Errorf("%w: %s, for tag: %s", ErrFailWrongRule, v.field, tag)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (sv *StringValidator) Len(arg string) error {
	condition, err := strconv.Atoi(arg)
	if err != nil {
		return ErrFailWrongLengthNumber
	}
	if condition != len(sv.value) {
		sv.v.addValidationError(fmt.Errorf("%w: need %d, have %d", ErrStringLength, condition, len(sv.value)))
	}
	return nil
}

func (sv *StringValidator) Regexp(arg string) error {
	r, err := regexp.Compile(arg)
	if err != nil {
		return ErrFailWrongRegexp
	}
	matched := r.MatchString(sv.value)
	if !matched {
		sv.v.addValidationError(fmt.Errorf("%w: %s", ErrStringNotMatchedRegexp, arg))
	}
	return nil
}

func (sv *StringValidator) In(arg string) error {
	if len(arg) == 0 {
		return ErrFailWrongSet
	}
	set := strings.Split(arg, ",")
	in := false
	for _, s := range set {
		if s == sv.value {
			in = true
			break
		}
	}
	if !in {
		sv.v.addValidationError(fmt.Errorf("%w: value %s, set %v", ErrStringNotIncludedSet, sv.value, set))
	}
	return nil
}
