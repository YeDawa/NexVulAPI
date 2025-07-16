package generator

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Slugify(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	text, _, _ = transform.String(t, text)

	re := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	text = re.ReplaceAllString(text, "")

	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, "-")

	text = strings.ToLower(text)
	text = strings.Trim(text, "-")

	re = regexp.MustCompile(`-+`)
	text = re.ReplaceAllString(text, "-")

	return text
}
