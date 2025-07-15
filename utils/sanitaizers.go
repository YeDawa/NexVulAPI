package utils

import "strings"

func SanitizeText(text string) string {
	replacer := strings.NewReplacer(
		"→", "->",
		"•", "-",
		"“", "\"",
		"”", "\"",
		"’", "'",
		"–", "-",
		"—", "--",
		"…", "...",
		"©", "(c)",
	)
	return replacer.Replace(text)
}
