package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/enescakir/emoji"

	"selectel/config"
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

func checkNoSpecialChars(msg string, cfg *config.Config) bool {
	if cfg == nil {
		cfg = config.Load("")
	}

	for _, r := range msg {
		// emojis
		if emoji.Exist(string(r)) {
			return false
		}

		// explicitly forbidden symbols from config
		if strings.ContainsRune(cfg.ForbiddenSymbols, r) {
			return false
		}
	}

	return true
}

func checkSensitive(msg string, cfg *config.Config) bool {
	if cfg == nil {
		cfg = config.Load("")
	}

	lower := strings.ToLower(msg)
	for _, kw := range cfg.SensitivePatterns {
		if strings.Contains(lower, kw) {
			return false
		}
	}
	return true
}
