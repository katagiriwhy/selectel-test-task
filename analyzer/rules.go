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

func checkEnglish(msg string) bool {
	for _, r := range msg {
		if unicode.In(r, unicode.Cyrillic) {
			return false
		}
	}
	return true
}
