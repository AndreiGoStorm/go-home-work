package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	UserProfile struct {
		ID       string   `json:"id" validate:"len:36"`
		UserInfo UserInfo `validate:"nested"`
		Identity int64    `validate:"min:1"`
	}

	Role struct {
		ID   int64  `validate:"min:1"`
		Name string `validate:"regexp:^\\w+$"`
		Ref  string
		UUID string `validate:"len:36"`
	}

	UserInfo struct {
		ID     int64  `validate:"min:1"`
		Role   Role   `validate:"nested"`
		Age    int    `validate:"min:18|max:50"`
		Phones string `validate:"len:11"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "b95b6a96-588a-4fcb-b169-0632a2c5bfc5",
				Name:   "Andrei",
				Age:    18,
				Email:  "andrei@example.com",
				Role:   "admin",
				Phones: []string{"57631726972", "42745657177"},
				meta:   nil,
			},
			nil,
		},
		{
			App{"2.3.4"},
			nil,
		},
		{
			Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "{body}",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

func TestNestedValidator(t *testing.T) {
	role := Role{
		ID:   325345634,
		Name: "Oleg",
		UUID: "b95b6a96-588a-4fcb-b169-0632a2c5bfc5",
	}

	userInfo := UserInfo{
		ID:     105121255,
		Role:   role,
		Age:    25,
		Phones: "57631726972",
	}

	userProfile := UserProfile{
		ID:       "2bb24361-d0ef-426f-91d2-be42210390ed",
		UserInfo: userInfo,
		Identity: 11254343348,
	}

	t.Run("nested validation", func(t *testing.T) {
		err := Validate(userProfile)
		require.NoError(t, err)
	})
}

func TestStructValidator(t *testing.T) {
	t.Run("interface is not a struct", func(t *testing.T) {
		tests := []struct {
			in interface{}
		}{
			{in: "string"},
			{in: 98},
			{in: 12.23},
			{in: []int{3, 4, 5}},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("struct case %d", i), func(t *testing.T) {
				err := Validate(tt.in)
				require.Error(t, err)
				require.ErrorIs(t, err, ErrFailNotStruct)
			})
		}
	})
}

type ValidatorTest struct {
	in            interface{}
	expectedErr   error
	expectedField string
}

func TestSimpleValidate(t *testing.T) {
	tests := getStringTestParams()
	tests = append(tests, getIntTestParams()...)
	tests = append(tests, getSliceTestParams()...)

	for i, tt := range tests {
		t.Run(fmt.Sprintf("simple case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			if tt.expectedErr != nil {
				var errs ValidationErrors
				if errors.As(err, &errs) {
					require.ErrorIs(t, errs[0].Err, tt.expectedErr)
					require.Equal(t, errs[0].Field, tt.expectedField)
				} else {
					require.Fail(t, "test failed")
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type (
	StringLen struct {
		Len string `validate:"len:7"`
	}

	StringRegexp struct {
		Regexp string `validate:"regexp:[0-9]+"`
	}

	StringIn struct {
		In string `validate:"in:foo,bar"`
	}

	IntMin struct {
		Min int `validate:"min:15"`
	}

	IntMax struct {
		Max int `validate:"max:50"`
	}

	IntIn struct {
		In int `validate:"in:1,2,3,4,5"`
	}

	StringSlice struct {
		Strings []string `validate:"len:5|in:first,third"`
	}

	IntSlice struct {
		Ints []int `validate:"min:1|max:50|in:5,10,15"`
	}
)

func getStringTestParams() []ValidatorTest {
	return []ValidatorTest{
		{
			in: StringLen{Len: "version"},
		},
		{
			in:            StringLen{Len: "app"},
			expectedErr:   ErrStringLength,
			expectedField: "Len",
		},
		{
			in: StringRegexp{Regexp: "23"},
		},
		{
			in:            StringRegexp{Regexp: "version"},
			expectedErr:   ErrStringNotMatchedRegexp,
			expectedField: "Regexp",
		},
		{
			in: StringIn{In: "foo"},
		},
		{
			in:            StringIn{In: "version"},
			expectedErr:   ErrStringNotIncludedSet,
			expectedField: "In",
		},
	}
}

func getIntTestParams() []ValidatorTest {
	return []ValidatorTest{
		{
			in: IntMin{Min: 15},
		},
		{
			in:            IntMin{Min: 10},
			expectedErr:   ErrIntLessThanMin,
			expectedField: "Min",
		},
		{
			in: IntMax{Max: 23},
		},
		{
			in:            IntMax{Max: 51},
			expectedErr:   ErrIntMoreThanMax,
			expectedField: "Max",
		},
		{
			in: IntIn{In: 2},
		},
		{
			in:            IntIn{In: 10},
			expectedErr:   ErrIntNotIncludedInSet,
			expectedField: "In",
		},
	}
}

func getSliceTestParams() []ValidatorTest {
	return []ValidatorTest{
		{
			in: StringSlice{Strings: []string{"first", "third"}},
		},
		{
			in:            StringSlice{Strings: []string{"first", "thitt"}},
			expectedErr:   ErrStringNotIncludedSet,
			expectedField: "Strings",
		},
		{
			in: IntSlice{Ints: []int{5, 10, 10, 15, 15, 15}},
		},
		{
			in:            IntSlice{Ints: []int{0, 5, 10, 15}},
			expectedErr:   ErrIntLessThanMin,
			expectedField: "Ints",
		},
	}
}

func TestFailValidate(t *testing.T) {
	tests := getStringFailParams()
	tests = append(tests, getIntFailParams()...)

	for i, tt := range tests {
		t.Run(fmt.Sprintf("fail validate case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

type (
	StringWrongRule struct {
		Wrong string `validate:"min:7"`
	}

	StringFailLen struct {
		Len string `validate:"len:W7"`
	}

	StringFailRegexp struct {
		Regexp string `validate:"regexp:stp("`
	}

	StringFailIn struct {
		In string `validate:"in:"`
	}

	IntWrongRule struct {
		Wrong int `validate:"len:7"`
	}

	IntFailMin struct {
		Min int `validate:"min:w"`
	}

	IntFailMax struct {
		Max int `validate:"max:r"`
	}

	IntFailIn struct {
		In int `validate:"in:slice"`
	}
)

func getStringFailParams() []ValidatorTest {
	return []ValidatorTest{
		{
			in:          StringWrongRule{Wrong: "version"},
			expectedErr: ErrFailWrongRule,
		},
		{
			in:          StringFailLen{Len: "version"},
			expectedErr: ErrFailWrongLengthNumber,
		},
		{
			in:          StringFailRegexp{Regexp: "version"},
			expectedErr: ErrFailWrongRegexp,
		},
		{
			in:          StringFailIn{In: "version"},
			expectedErr: ErrFailWrongSet,
		},
	}
}

func getIntFailParams() []ValidatorTest {
	return []ValidatorTest{
		{
			in:          IntWrongRule{Wrong: 100},
			expectedErr: ErrFailWrongRule,
		},
		{
			in:          IntFailMin{Min: 100},
			expectedErr: ErrFailWrongMinNumber,
		},
		{
			in:          IntFailMax{Max: 100},
			expectedErr: ErrFailWrongMaxNumber,
		},
		{
			in:          IntFailIn{In: 100},
			expectedErr: ErrFailSetNumber,
		},
	}
}
