package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrStringLengthNumber     = errors.New("value for length is not a number")
	ErrStringLength           = errors.New("length of string is wrong")
	ErrStringRegexp           = errors.New("regexp is wrong")
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

func (sv *StringValidator) validate(v *Validator) {
	sv.prepareParams(v)
	for _, tag := range v.tags {
		args := v.parseTag(tag)
		if args == nil {
			return
		}

		switch args[0] {
		case "len":
			sv.Len(args[1])
		case "regexp":
			sv.Regexp(args[1])
		case "in":
			sv.In(args[1])
		default:
			sv.v.addValidationError(fmt.Errorf("%w: %s, for tag: %s", ErrUndefinedRule, v.field, tag))
		}
	}
}

func (sv *StringValidator) Len(arg string) {
	condition, err := strconv.Atoi(arg)
	if err != nil {
		sv.v.addValidationError(ErrStringLengthNumber)
		return
	}
	if condition != len(sv.value) {
		sv.v.addValidationError(fmt.Errorf("%w: need %d, have %d", ErrStringLength, condition, len(sv.value)))
	}
}

func (sv *StringValidator) Regexp(arg string) {
	r, err := regexp.Compile(arg)
	if err != nil {
		sv.v.addValidationError(ErrStringRegexp)
		return
	}
	matched := r.MatchString(sv.value)
	if !matched {
		sv.v.addValidationError(fmt.Errorf("%w: %s", ErrStringNotMatchedRegexp, arg))
	}
}

func (sv *StringValidator) In(arg string) {
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
}
