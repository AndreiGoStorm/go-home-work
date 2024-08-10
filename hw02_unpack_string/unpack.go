package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const slash = '\\'

type Token struct {
	token       string
	unpacked    string
	isForUnpack bool
}

func Unpack(s string) (string, error) {
	tokens, err := tokenize([]rune(s))
	if err != nil {
		return "", err
	}
	err = unpackTokens(tokens)
	if err != nil {
		return "", err
	}
	return unpackString(tokens), nil
}

func tokenize(runes []rune) ([]Token, error) {
	tokens := make([]Token, 0, len(runes))
	for i := 0; i < len(runes); i++ {
		if unicode.IsDigit(runes[i]) {
			if i != 0 {
				tokens = append(tokens, Token{token: string(runes[i]), isForUnpack: true})
				continue
			}
			return nil, ErrInvalidString
		}
		if runes[i] == slash {
			if i == len(runes)-1 {
				return nil, ErrInvalidString
			}
			i++
		}
		tokens = append(tokens, Token{token: string(runes[i])})
	}
	return tokens, nil
}

func unpackTokens(tokens []Token) error {
	for index, value := range tokens {
		tokens[index].unpacked = value.token
		if value.isForUnpack {
			prevToken := &tokens[index-1]
			if prevToken.isForUnpack {
				return ErrInvalidString
			}
			length, _ := strconv.Atoi(value.token)
			tokens[index].unpacked = strings.Repeat(prevToken.unpacked, length)
			prevToken.unpacked = ""
		}
	}
	return nil
}

func unpackString(tokens []Token) string {
	var result string
	for _, value := range tokens {
		result += value.unpacked
	}
	return result
}
