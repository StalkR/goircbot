package url

import (
	"errors"
	"html"
	"regexp"
)

var defaultRE = regexp.MustCompile(`(?i)<title[^>]*>([^<]+)<`)

func handleDefault(url string) (string, error) {
	body, err := get(url)
	if err != nil {
		return "", err
	}
	text := defaultRE.FindStringSubmatch(body)
	if text == nil {
		return "", errors.New("url: cannot parse title")
	}
	return trim(html.UnescapeString(text[1])), nil
}
