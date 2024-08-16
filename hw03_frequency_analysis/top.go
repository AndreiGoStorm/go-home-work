package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Word struct {
	count int
	value string
}

const maxCount = 10

var regExp = regexp.MustCompile(`(?i)([а-яa-z\-\,\.]*[а-яa-z])|(-{2,})`)

func Top10(text string) []string {
	fields := regExp.FindAllString(text, -1)
	wordsMap := make(map[string]int, len(fields))
	for _, value := range fields {
		valueLower := strings.ToLower(value)
		wordsMap[valueLower]++
	}
	words := make([]Word, 0, len(wordsMap))
	for key, count := range wordsMap {
		words = append(words, Word{count, key})
	}

	return filterTop10(&words)
}

func filterTop10(words *[]Word) []string {
	sort.Slice(*words, func(i, j int) bool {
		if (*words)[i].count == (*words)[j].count {
			return (*words)[i].value < (*words)[j].value
		}
		return (*words)[i].count > (*words)[j].count
	})

	length := len(*words)
	if length > maxCount {
		length = maxCount
	}
	top10 := make([]string, 0, length)
	for i := 0; i < length; i++ {
		top10 = append(top10, (*words)[i].value)
	}

	return top10
}
