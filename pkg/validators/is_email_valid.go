package validators

import (
	"regexp"
	"unicode/utf8"
)

func IsEmailValid(email string) bool {
	length := utf8.RuneCountInString(email)
	pattern := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if length < 5 || length > 64 || !pattern.MatchString(email) {
		return false
	}

	return true
}
