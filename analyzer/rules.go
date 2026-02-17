package analyzer

import (
	"unicode"
	"unicode/utf8"
)

func checkLowercase(msg string) bool {
	if msg == "" {
		return true
	}

	r, _ := utf8.DecodeRuneInString(msg)
	return unicode.IsLower(r)
}
