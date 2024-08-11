package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "abccd0", expected: "abcc"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "a5b2ctd3", expected: "aaaaabbctddd"},
		{input: "d4ar0b7dwa2a1a0", expected: "ddddabbbbbbbdwaaa"},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `qwe\4\53`, expected: `qwe4555`},
		{input: `qwe\\\3\\`, expected: `qwe\3\`},
		{input: `q\24ar\30\22b1d2a1a\0`, expected: `q2222ar22bddaa0`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `\\\`, `r\\\42\`, `\a`, `\`, `\2\b`, `d\re\5p`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
