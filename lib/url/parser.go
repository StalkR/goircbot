package url

import (
	"errors"
	"html"
	"regexp"
	"strings"
)

// Parsers contains an ordered list of supported url title Parsers.
var Parsers = []Parser{
	&twitter{},
}

// A Parser can Parse() a body to extract a title if a given URL Match()es.
type Parser interface {
	Match(url string) bool
	Parse(body string) (string, error)
}

var titleRE = regexp.MustCompile(`<title[^>]*>([^<]+)<`)

// ParseTitle parses a title from an URL and content.
func ParseTitle(url, body string) (string, error) {
	for _, p := range Parsers {
		if p.Match(url) {
			return p.Parse(body)
		}
	}
	// If no match, default parser applies.
	matches := titleRE.FindStringSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("url: cannot parse title")
	}
	return Trim(html.UnescapeString(matches[1])), nil
}

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
