package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/enescakir/emoji"
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

func checkNoSpecialChars(msg string) bool {
	for _, r := range msg {
		if unicode.IsSymbol(r) || unicode.IsPunct(r) {
			if r != ' ' || emoji.Exist(string(r)) {
				return false
			}
		}
	}
	return true
}

var sensitiveKeywords = []string{
	"password",
	"token",
	"api_key",
	"secret",
}

func checkSensitive(msg string) bool {
	lower := strings.ToLower(msg)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return false
		}
	}
	return true
}
