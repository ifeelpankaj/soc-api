package logger

import (
	"regexp"
	"strings"
	"unicode"
)

var credentialPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(bearer|basic)\s+[^\s,;]+`),
	regexp.MustCompile(`(?i)"?(password|passwd|secret|token|api[_-]?key|authorization|cookie)"?\s*[=:]\s*("[^"]*"|'[^']*'|[^\s,;]+)`),
	regexp.MustCompile(`(?i)://[^/@\s]+:[^/@\s]+@`),
}

// Sanitize bounds diagnostic strings and removes common credential formats.
// Callers must still avoid passing bodies, headers, or arbitrary payloads.
func Sanitize(value string) string {
	for _, pattern := range credentialPatterns {
		value = pattern.ReplaceAllString(value, "[REDACTED]")
	}
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, value)
	runes := []rune(value)
	if len(runes) > 2048 {
		value = string(runes[:2048])
	}
	return value
}
