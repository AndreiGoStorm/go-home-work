package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	reg, err := regexp.Compile("(?i).+@(.+\\." + domain + ")")
	if err != nil {
		return nil, err
	}
	return countDomains(r, reg)
}

func countDomains(r io.Reader, reg *regexp.Regexp) (DomainStat, error) {
	rdr := bufio.NewReader(r)

	result := make(DomainStat, 100)
	var user User
	for {
		line, _, err := rdr.ReadLine()
		if err == io.EOF { //nolint:errorlint
			break
		}
		if err != nil {
			return nil, err
		}

		err = user.UnmarshalJSON(line)
		if err != nil {
			return nil, err
		}

		matched := reg.FindAllStringSubmatch(user.Email, 2)
		if matched == nil {
			continue
		}
		emailDomain := strings.ToLower(matched[0][1])
		result[emailDomain]++
	}
	return result, nil
}
