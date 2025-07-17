package utils

import (
	"encoding/json"
	"strings"
)

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

func ToJSONString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJSONString(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}
