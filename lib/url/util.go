package url

import (
	"regexp"
	"strings"
)

var whitespaceRE = regexp.MustCompile(`\s+`)

// Trim removes all white spaces including duplicates in a string.
func Trim(s string) string {
	return strings.TrimSpace(whitespaceRE.ReplaceAllString(s, " "))
}

var htmlTagsRE = regexp.MustCompile(`<[^>]*>`)

// StripTags removes all HTML tags in a string.
func StripTags(s string) string {
	return htmlTagsRE.ReplaceAllString(s, "")
}
